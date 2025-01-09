import { useQuery } from '@tanstack/react-query';
import Schedule from '@/model/schedule.ts';

export default function useFetchSchedule(scheduleId?: string) {
  return useQuery({
    queryKey: ['schedule', scheduleId],
    queryFn: () => fetchSchedule(scheduleId),
    retry: false,
  });
}

const fetchSchedule = async (scheduleId?: string): Promise<Schedule | null> => {
  if (!scheduleId) return null;
  const response = await fetch(
    `${import.meta.env.VITE_API_URL}/schedules/${scheduleId}`
  );
  return response.json();
};
