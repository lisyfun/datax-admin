import { ref } from 'vue';

export default function useLoading(initialValue = false) {
  const loading = ref(initialValue);

  const setLoading = (value: boolean) => {
    loading.value = value;
  };

  return {
    loading,
    setLoading,
  };
}
