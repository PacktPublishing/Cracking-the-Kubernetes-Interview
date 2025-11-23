package v1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// BookSpec defines the desired state of Book
type BookSpec struct {
    // Book title (required, non-empty string)
    // +kubebuilder:validation:Required
    // +kubebuilder:validation:MinLength=1
    // +kubebuilder:validation:MaxLength=200
    Book string `json:"book"`

    // Publication year (required, between 1900-2100)
    // +kubebuilder:validation:Required
    // +kubebuilder:validation:Minimum=1900
    // +kubebuilder:validation:Maximum=2100
    Year int `json:"year"`
}

// BookStatus defines the observed state of Book
type BookStatus struct {
    // INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
    // Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Book is the Schema for the books API
type Book struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   BookSpec   `json:"spec,omitempty"`
    Status BookStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BookList contains a list of Book
type BookList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []Book `json:"items"`
}

func init() {
    SchemeBuilder.Register(&Book{}, &BookList{})
}
