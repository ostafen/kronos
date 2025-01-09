import { useMutation } from '@tanstack/react-query';

export default function useTriggerSchedule() {
  return useMutation({
    mutationKey: ['triggerSchedule'],
    mutationFn: (scheduleId: string) =>
      fetch(`${import.meta.env.VITE_API_URL}/schedules/${scheduleId}/trigger`, {
        method: 'POST',
      }),
  });
}
