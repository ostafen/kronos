export default interface Schedule {
  id: string;
  title: string;
  status: ScheduleStatus;
  description: string;
  cronExpr?: string;
  url: string;
  metadata: string;
  isRecurring: boolean;
  createdAt: string;
  runAt: string;
  startAt: string;
  endAt: string;
}

export type ScheduleStatus = 'not_started' | 'active' | 'paused' | 'expired';

export type NewSchedule = Omit<
  Partial<Schedule>,
  'id' | 'status' | 'createdAt'
>;
