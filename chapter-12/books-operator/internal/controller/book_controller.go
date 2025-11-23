package controllers

import (
    "context"
    "fmt"

    "k8s.io/apimachinery/pkg/api/resource"
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/log"

    packtv1 "github.com/PacktPublishing/Kubernetes-Interview-Guide/chapter-12/books-operator/api/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BookReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=packt.com,resources=books,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=packt.com,resources=books/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=packt.com,resources=books/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main Kubernetes reconciliation loop
func (r *BookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := log.FromContext(ctx)

    // Fetch the Book instance
    book := &packtv1.Book{}
    err := r.Get(ctx, req.NamespacedName, book)
    if err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    // Log the reconciliation
    log.Info("Reconciling Book", "name", book.Name, "namespace", book.Namespace, "book", book.Spec.Book, "year", book.Spec.Year)

    // Define a new Pod object
    pod := &corev1.Pod{
        ObjectMeta: metav1.ObjectMeta{
            Name:      book.Name + "-pod",
            Namespace: book.Namespace, // Use the Book CR's namespace
            OwnerReferences: []metav1.OwnerReference{
                *metav1.NewControllerRef(book, packtv1.GroupVersion.WithKind("Book")),
            },
        },
        Spec: corev1.PodSpec{
            Containers: []corev1.Container{
                {
                    Name:  "busybox",
                    Image: "busybox:1.36", // Pin image version for reproducibility
                    Command: []string{
                        "sh",
                        "-c",
                        fmt.Sprintf("while true; do echo Book: %s, Year: %d; sleep 1; done", book.Spec.Book, book.Spec.Year),
                    },
                    Resources: corev1.ResourceRequirements{
                        Requests: corev1.ResourceList{
                            corev1.ResourceCPU:    resource.MustParse("100m"),
                            corev1.ResourceMemory: resource.MustParse("128Mi"),
                        },
                        Limits: corev1.ResourceList{
                            corev1.ResourceCPU:    resource.MustParse("200m"),
                            corev1.ResourceMemory: resource.MustParse("256Mi"),
                        },
                    },
                },
            },
        },
    }

    // Check if the Pod already exists
    found := &corev1.Pod{}
    err = r.Get(ctx, client.ObjectKey{Name: pod.Name, Namespace: pod.Namespace}, found)
    if err != nil && client.IgnoreNotFound(err) != nil {
        log.Error(err, "Failed to get Pod")
        return ctrl.Result{}, err
    }

    if err == nil {
        // Pod already exists - don't requeue
        log.Info("Pod already exists", "pod", pod.Name, "namespace", pod.Namespace)
        return ctrl.Result{}, nil
    }

    // Create the Pod
    log.Info("Creating Pod", "pod", pod.Name, "namespace", pod.Namespace)
    err = r.Create(ctx, pod)
    if err != nil {
        log.Error(err, "Failed to create Pod", "pod", pod.Name, "namespace", pod.Namespace)
        return ctrl.Result{}, err
    }

    log.Info("Pod created successfully", "pod", pod.Name, "namespace", pod.Namespace)
    // Pod created successfully - don't requeue
    return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BookReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&packtv1.Book{}).
        Owns(&corev1.Pod{}). // Watch pods owned by Book CRs
        Complete(r)
}
