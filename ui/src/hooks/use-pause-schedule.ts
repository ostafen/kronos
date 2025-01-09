import { useMutation } from '@tanstack/react-query';

export default function usePauseSchedule() {
  return useMutation({
    mutationKey: ['pauseSchedule'],
    mutationFn: (scheduleId: string) =>
      fetch(`${import.meta.env.VITE_API_URL}/schedules/${scheduleId}/pause`, {
        method: 'POST',
      }),
  });
}
