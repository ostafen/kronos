import {useQuery} from "@tanstack/react-query";
import Schedule from "@/model/schedule.ts";

export default function useFetchSchedules() {
    return useQuery({
        queryKey: ["schedules"],
        queryFn: fetchSchedules,
    });
}

async function fetchSchedules(): Promise<Schedule[]> {
    try {
        const response = await fetch('http://localhost:9175/schedules');

        if (!response.ok) {
            console.error("Error fetching schedules");
            return [];
        }

        return response.json();
    } catch (error) {
        if (error instanceof Error) {
            console.error(error.message);
        }

        return [];
    }
}