import { useMutation } from '@tanstack/react-query';

export default function useDeleteSchedule() {
  return useMutation({
    mutationKey: ['deleteSchedule'],
    mutationFn: (id: string) =>
      fetch(`${import.meta.env.VITE_API_URL}/schedules/${id}`, {
        method: 'DELETE',
      }),
  });
}
