import { useQueryClient } from '@tanstack/react-query';

export default function useInvalidateSchedules(
  callback: (...params: unknown[]) => Promise<unknown>
) {
  const queryClient = useQueryClient();

  return async (...params: unknown[]) => {
    await callback(params);
    await queryClient.invalidateQueries({ queryKey: ['schedules'] });
  };
}
