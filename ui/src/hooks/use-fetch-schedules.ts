import { useQuery } from "@tanstack/react-query";
import Schedule from "@/model/schedule.ts";

export default function useFetchSchedules() {
  return useQuery({
    queryKey: ["schedules"],
    queryFn: fetchSchedules,
  });
}

async function fetchSchedules(): Promise<Schedule[]> {
  const response = await fetch(`${import.meta.env.VITE_API_URL}/schedules`);
  return response.json();
}
