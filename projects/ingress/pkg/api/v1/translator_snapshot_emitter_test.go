// Code generated by solo-kit. DO NOT EDIT.

// +build solokit

package v1

import (
	"context"
	"os"
	"time"

	gloo_solo_io "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	kuberc "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/setup"
	"k8s.io/client-go/rest"

	// Needed to run tests in GKE
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	// From https://github.com/kubernetes/client-go/blob/53c7adfd0294caa142d961e1f780f74081d5b15f/examples/out-of-cluster-client-configuration/main.go#L31
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/cache"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var _ = Describe("V1Emitter", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		namespace1     string
		namespace2     string
		name1, name2   = "angela" + helpers.RandString(3), "bob" + helpers.RandString(3)
		cfg            *rest.Config
		emitter        TranslatorEmitter
		secretClient   gloo_solo_io.SecretClient
		upstreamClient gloo_solo_io.UpstreamClient
		ingressClient  IngressClient
	)

	BeforeEach(func() {
		namespace1 = helpers.RandString(8)
		namespace2 = helpers.RandString(8)
		var err error
		cfg, err = kubeutils.GetConfig("", "")
		Expect(err).NotTo(HaveOccurred())
		err = setup.SetupKubeForTest(namespace1)
		Expect(err).NotTo(HaveOccurred())
		err = setup.SetupKubeForTest(namespace2)
		Expect(err).NotTo(HaveOccurred())
		var kube kubernetes.Interface
		// Secret Constructor
		kube, err = kubernetes.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())

		kcache, err := cache.NewKubeCoreCache(context.TODO(), kube)
		Expect(err).NotTo(HaveOccurred())
		secretClientFactory := &factory.KubeSecretClientFactory{
			Clientset: kube,
			Cache:     kcache,
		}
		secretClient, err = gloo_solo_io.NewSecretClient(secretClientFactory)
		Expect(err).NotTo(HaveOccurred())
		// Upstream Constructor
		upstreamClientFactory := &factory.KubeResourceClientFactory{
			Crd:         gloo_solo_io.UpstreamCrd,
			Cfg:         cfg,
			SharedCache: kuberc.NewKubeCache(context.TODO()),
		}
		upstreamClient, err = gloo_solo_io.NewUpstreamClient(upstreamClientFactory)
		Expect(err).NotTo(HaveOccurred())
		// Ingress Constructor
		kube, err = kubernetes.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())

		kcache, err := cache.NewKubeCoreCache(context.TODO(), kube)
		Expect(err).NotTo(HaveOccurred())
		ingressClientFactory := &factory.KubeConfigMapClientFactory{
			Clientset: kube,
			Cache:     kcache,
		}
		ingressClient, err = NewIngressClient(ingressClientFactory)
		Expect(err).NotTo(HaveOccurred())
		emitter = NewTranslatorEmitter(secretClient, upstreamClient, ingressClient)
	})
	AfterEach(func() {
		setup.TeardownKube(namespace1)
		setup.TeardownKube(namespace2)
	})
	It("tracks snapshots on changes to any resource", func() {
		ctx := context.Background()
		err := emitter.Register()
		Expect(err).NotTo(HaveOccurred())

		snapshots, errs, err := emitter.Snapshots([]string{namespace1, namespace2}, clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Second,
		})
		Expect(err).NotTo(HaveOccurred())

		var snap *TranslatorSnapshot

		/*
			Secret
		*/

		assertSnapshotSecrets := func(expectSecrets gloo_solo_io.SecretList, unexpectSecrets gloo_solo_io.SecretList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectSecrets {
						if _, err := snap.Secrets.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectSecrets {
						if _, err := snap.Secrets.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := secretClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := secretClient.List(namespace2, clients.ListOpts{})
					combined := gloo_solo_io.SecretsByNamespace{
						namespace1: nsList1,
						namespace2: nsList2,
					}
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		secret1a, err := secretClient.Write(gloo_solo_io.NewSecret(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		secret1b, err := secretClient.Write(gloo_solo_io.NewSecret(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotSecrets(gloo_solo_io.SecretList{secret1a, secret1b}, nil)
		secret2a, err := secretClient.Write(gloo_solo_io.NewSecret(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		secret2b, err := secretClient.Write(gloo_solo_io.NewSecret(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotSecrets(gloo_solo_io.SecretList{secret1a, secret1b, secret2a, secret2b}, nil)

		err = secretClient.Delete(secret2a.Metadata.Namespace, secret2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = secretClient.Delete(secret2b.Metadata.Namespace, secret2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotSecrets(gloo_solo_io.SecretList{secret1a, secret1b}, gloo_solo_io.SecretList{secret2a, secret2b})

		err = secretClient.Delete(secret1a.Metadata.Namespace, secret1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = secretClient.Delete(secret1b.Metadata.Namespace, secret1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotSecrets(nil, gloo_solo_io.SecretList{secret1a, secret1b, secret2a, secret2b})

		/*
			Upstream
		*/

		assertSnapshotUpstreams := func(expectUpstreams gloo_solo_io.UpstreamList, unexpectUpstreams gloo_solo_io.UpstreamList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectUpstreams {
						if _, err := snap.Upstreams.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectUpstreams {
						if _, err := snap.Upstreams.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := upstreamClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := upstreamClient.List(namespace2, clients.ListOpts{})
					combined := gloo_solo_io.UpstreamsByNamespace{
						namespace1: nsList1,
						namespace2: nsList2,
					}
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		upstream1a, err := upstreamClient.Write(gloo_solo_io.NewUpstream(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		upstream1b, err := upstreamClient.Write(gloo_solo_io.NewUpstream(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotUpstreams(gloo_solo_io.UpstreamList{upstream1a, upstream1b}, nil)
		upstream2a, err := upstreamClient.Write(gloo_solo_io.NewUpstream(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		upstream2b, err := upstreamClient.Write(gloo_solo_io.NewUpstream(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotUpstreams(gloo_solo_io.UpstreamList{upstream1a, upstream1b, upstream2a, upstream2b}, nil)

		err = upstreamClient.Delete(upstream2a.Metadata.Namespace, upstream2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = upstreamClient.Delete(upstream2b.Metadata.Namespace, upstream2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotUpstreams(gloo_solo_io.UpstreamList{upstream1a, upstream1b}, gloo_solo_io.UpstreamList{upstream2a, upstream2b})

		err = upstreamClient.Delete(upstream1a.Metadata.Namespace, upstream1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = upstreamClient.Delete(upstream1b.Metadata.Namespace, upstream1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotUpstreams(nil, gloo_solo_io.UpstreamList{upstream1a, upstream1b, upstream2a, upstream2b})

		/*
			Ingress
		*/

		assertSnapshotIngresses := func(expectIngresses IngressList, unexpectIngresses IngressList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectIngresses {
						if _, err := snap.Ingresses.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectIngresses {
						if _, err := snap.Ingresses.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := ingressClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := ingressClient.List(namespace2, clients.ListOpts{})
					combined := IngressesByNamespace{
						namespace1: nsList1,
						namespace2: nsList2,
					}
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		ingress1a, err := ingressClient.Write(NewIngress(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		ingress1b, err := ingressClient.Write(NewIngress(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotIngresses(IngressList{ingress1a, ingress1b}, nil)
		ingress2a, err := ingressClient.Write(NewIngress(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		ingress2b, err := ingressClient.Write(NewIngress(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotIngresses(IngressList{ingress1a, ingress1b, ingress2a, ingress2b}, nil)

		err = ingressClient.Delete(ingress2a.Metadata.Namespace, ingress2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = ingressClient.Delete(ingress2b.Metadata.Namespace, ingress2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotIngresses(IngressList{ingress1a, ingress1b}, IngressList{ingress2a, ingress2b})

		err = ingressClient.Delete(ingress1a.Metadata.Namespace, ingress1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = ingressClient.Delete(ingress1b.Metadata.Namespace, ingress1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotIngresses(nil, IngressList{ingress1a, ingress1b, ingress2a, ingress2b})
	})
	It("tracks snapshots on changes to any resource using AllNamespace", func() {
		ctx := context.Background()
		err := emitter.Register()
		Expect(err).NotTo(HaveOccurred())

		snapshots, errs, err := emitter.Snapshots([]string{""}, clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Second,
		})
		Expect(err).NotTo(HaveOccurred())

		var snap *TranslatorSnapshot

		/*
			Secret
		*/

		assertSnapshotSecrets := func(expectSecrets gloo_solo_io.SecretList, unexpectSecrets gloo_solo_io.SecretList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectSecrets {
						if _, err := snap.Secrets.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectSecrets {
						if _, err := snap.Secrets.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := secretClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := secretClient.List(namespace2, clients.ListOpts{})
					combined := gloo_solo_io.SecretsByNamespace{
						namespace1: nsList1,
						namespace2: nsList2,
					}
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		secret1a, err := secretClient.Write(gloo_solo_io.NewSecret(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		secret1b, err := secretClient.Write(gloo_solo_io.NewSecret(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotSecrets(gloo_solo_io.SecretList{secret1a, secret1b}, nil)
		secret2a, err := secretClient.Write(gloo_solo_io.NewSecret(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		secret2b, err := secretClient.Write(gloo_solo_io.NewSecret(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotSecrets(gloo_solo_io.SecretList{secret1a, secret1b, secret2a, secret2b}, nil)

		err = secretClient.Delete(secret2a.Metadata.Namespace, secret2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = secretClient.Delete(secret2b.Metadata.Namespace, secret2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotSecrets(gloo_solo_io.SecretList{secret1a, secret1b}, gloo_solo_io.SecretList{secret2a, secret2b})

		err = secretClient.Delete(secret1a.Metadata.Namespace, secret1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = secretClient.Delete(secret1b.Metadata.Namespace, secret1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotSecrets(nil, gloo_solo_io.SecretList{secret1a, secret1b, secret2a, secret2b})

		/*
			Upstream
		*/

		assertSnapshotUpstreams := func(expectUpstreams gloo_solo_io.UpstreamList, unexpectUpstreams gloo_solo_io.UpstreamList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectUpstreams {
						if _, err := snap.Upstreams.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectUpstreams {
						if _, err := snap.Upstreams.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := upstreamClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := upstreamClient.List(namespace2, clients.ListOpts{})
					combined := gloo_solo_io.UpstreamsByNamespace{
						namespace1: nsList1,
						namespace2: nsList2,
					}
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		upstream1a, err := upstreamClient.Write(gloo_solo_io.NewUpstream(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		upstream1b, err := upstreamClient.Write(gloo_solo_io.NewUpstream(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotUpstreams(gloo_solo_io.UpstreamList{upstream1a, upstream1b}, nil)
		upstream2a, err := upstreamClient.Write(gloo_solo_io.NewUpstream(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		upstream2b, err := upstreamClient.Write(gloo_solo_io.NewUpstream(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotUpstreams(gloo_solo_io.UpstreamList{upstream1a, upstream1b, upstream2a, upstream2b}, nil)

		err = upstreamClient.Delete(upstream2a.Metadata.Namespace, upstream2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = upstreamClient.Delete(upstream2b.Metadata.Namespace, upstream2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotUpstreams(gloo_solo_io.UpstreamList{upstream1a, upstream1b}, gloo_solo_io.UpstreamList{upstream2a, upstream2b})

		err = upstreamClient.Delete(upstream1a.Metadata.Namespace, upstream1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = upstreamClient.Delete(upstream1b.Metadata.Namespace, upstream1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotUpstreams(nil, gloo_solo_io.UpstreamList{upstream1a, upstream1b, upstream2a, upstream2b})

		/*
			Ingress
		*/

		assertSnapshotIngresses := func(expectIngresses IngressList, unexpectIngresses IngressList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectIngresses {
						if _, err := snap.Ingresses.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectIngresses {
						if _, err := snap.Ingresses.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := ingressClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := ingressClient.List(namespace2, clients.ListOpts{})
					combined := IngressesByNamespace{
						namespace1: nsList1,
						namespace2: nsList2,
					}
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		ingress1a, err := ingressClient.Write(NewIngress(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		ingress1b, err := ingressClient.Write(NewIngress(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotIngresses(IngressList{ingress1a, ingress1b}, nil)
		ingress2a, err := ingressClient.Write(NewIngress(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		ingress2b, err := ingressClient.Write(NewIngress(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotIngresses(IngressList{ingress1a, ingress1b, ingress2a, ingress2b}, nil)

		err = ingressClient.Delete(ingress2a.Metadata.Namespace, ingress2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = ingressClient.Delete(ingress2b.Metadata.Namespace, ingress2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotIngresses(IngressList{ingress1a, ingress1b}, IngressList{ingress2a, ingress2b})

		err = ingressClient.Delete(ingress1a.Metadata.Namespace, ingress1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = ingressClient.Delete(ingress1b.Metadata.Namespace, ingress1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotIngresses(nil, IngressList{ingress1a, ingress1b, ingress2a, ingress2b})
	})
})
