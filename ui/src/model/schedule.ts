export default interface Schedule {
    id: string;
    title: string;
    status: ScheduleStatus;
    description: string;
    cronExpr: string;
    url: string;
    metadata: unknown;
    isRecurring: boolean;
    createdAt: string;
    runAt: string;
    startAt: string;
    endAt: string;
}

export type ScheduleStatus = 'not_started' | 'active' | 'paused' | 'expired'