import { useMutation } from '@tanstack/react-query';

export default function useResumeSchedule() {
  return useMutation({
    mutationKey: ['resumeSchedule'],
    mutationFn: (scheduleId: string) =>
      fetch(`${import.meta.env.VITE_API_URL}/schedules/${scheduleId}/resume`, {
        method: 'POST',
      }),
  });
}
