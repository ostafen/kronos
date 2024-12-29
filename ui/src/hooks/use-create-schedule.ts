import { NewSchedule } from "@/model/schedule";
import { useMutation } from "@tanstack/react-query";

export const createSchedule = async (schedule: NewSchedule) => {
  const headers = new Headers();
  headers.append("Content-Type", "application/json");

  return fetch(`${import.meta.env.VITE_API_URL}/schedules`, {
    method: "POST",
    body: JSON.stringify(schedule),
    headers,
  });
};

export default function useCreateSchedule() {
  return useMutation<Response, unknown, NewSchedule>({
    mutationKey: ["createSchedule"],
    mutationFn: (schedule) => createSchedule(schedule),
  });
}
