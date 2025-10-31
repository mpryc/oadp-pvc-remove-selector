package plugin

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
)

func TestAppliesTo(t *testing.T) {
	plugin := NewPVCRemoveSelectorPlugin(logrus.New())
	selector, err := plugin.AppliesTo()

	require.NoError(t, err)
	assert.Equal(t, []string{"persistentvolumeclaims"}, selector.IncludedResources)
}

func TestExecute_RemovesSelector(t *testing.T) {
	plugin := NewPVCRemoveSelectorPlugin(logrus.New())

	// Create a PVC with a selector
	pvc := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "PersistentVolumeClaim",
			"metadata": map[string]interface{}{
				"name":      "test-pvc",
				"namespace": "test-namespace",
			},
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"environment": "production",
					},
				},
				"accessModes": []interface{}{"ReadWriteOnce"},
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"storage": "10Gi",
					},
				},
			},
			"status": map[string]interface{}{
				"phase": "Bound",
			},
		},
	}

	input := &velero.RestoreItemActionExecuteInput{
		Item: pvc,
	}

	output, err := plugin.Execute(input)
	require.NoError(t, err)
	require.NotNil(t, output)
	require.NotNil(t, output.UpdatedItem)

	// Convert to *unstructured.Unstructured to access Object field
	updatedItem, ok := output.UpdatedItem.(*unstructured.Unstructured)
	require.True(t, ok, "UpdatedItem should be *unstructured.Unstructured")

	// Verify selector was removed
	_, found, err := unstructured.NestedMap(updatedItem.Object, "spec", "selector")
	require.NoError(t, err)
	assert.False(t, found, "selector should be removed")

	// Verify status was removed
	_, found, err = unstructured.NestedMap(updatedItem.Object, "status")
	require.NoError(t, err)
	assert.False(t, found, "status should be removed")

	// Verify other fields are preserved
	accessModes, found, err := unstructured.NestedSlice(updatedItem.Object, "spec", "accessModes")
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, []interface{}{"ReadWriteOnce"}, accessModes)
}

func TestExecute_NoSelector(t *testing.T) {
	plugin := NewPVCRemoveSelectorPlugin(logrus.New())

	// Create a PVC without a selector
	pvc := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "PersistentVolumeClaim",
			"metadata": map[string]interface{}{
				"name":      "test-pvc",
				"namespace": "test-namespace",
			},
			"spec": map[string]interface{}{
				"accessModes": []interface{}{"ReadWriteOnce"},
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"storage": "10Gi",
					},
				},
			},
		},
	}

	input := &velero.RestoreItemActionExecuteInput{
		Item: pvc,
	}

	output, err := plugin.Execute(input)
	require.NoError(t, err)
	require.NotNil(t, output)
	require.NotNil(t, output.UpdatedItem)

	// Convert to *unstructured.Unstructured to access Object field
	updatedItem, ok := output.UpdatedItem.(*unstructured.Unstructured)
	require.True(t, ok, "UpdatedItem should be *unstructured.Unstructured")

	// Verify the PVC is processed successfully even without a selector
	accessModes, found, err := unstructured.NestedSlice(updatedItem.Object, "spec", "accessModes")
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, []interface{}{"ReadWriteOnce"}, accessModes)
}

func TestExecute_DoesNotMutateInput(t *testing.T) {
	plugin := NewPVCRemoveSelectorPlugin(logrus.New())

	// Create a PVC with a selector
	pvc := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "PersistentVolumeClaim",
			"metadata": map[string]interface{}{
				"name":      "test-pvc",
				"namespace": "test-namespace",
			},
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"environment": "production",
					},
				},
			},
		},
	}

	input := &velero.RestoreItemActionExecuteInput{
		Item: pvc,
	}

	_, err := plugin.Execute(input)
	require.NoError(t, err)

	// Convert to *unstructured.Unstructured to access Object field
	inputItem, ok := input.Item.(*unstructured.Unstructured)
	require.True(t, ok, "input.Item should be *unstructured.Unstructured")

	// Verify original input was not mutated
	selector, found, err := unstructured.NestedMap(inputItem.Object, "spec", "selector")
	require.NoError(t, err)
	assert.True(t, found, "original input should still have selector")
	assert.NotNil(t, selector)
}
