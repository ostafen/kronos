import Schedule, {ScheduleStatus} from '@/model/schedule';
import {Badge, ColorPalette, Table as ChakraTable} from '@chakra-ui/react';
import dayjs from "dayjs";
import {ReactNode} from "react";

interface SchedulesTableProps {
    schedules: Schedule[];
}

export default function SchedulesTable(props: SchedulesTableProps) {
    const {schedules} = props;


    return (
        <ChakraTable.Root>
            <ChakraTable.Header>
                <ChakraTable.Row>
                    {scheduleFields.map((field) => (
                        <ChakraTable.ColumnHeader key={field.key}>{field.value}</ChakraTable.ColumnHeader>
                    ))}
                </ChakraTable.Row>
            </ChakraTable.Header>
            <ChakraTable.Body>
                {schedules.map((schedule) => (
                    <ChakraTable.Row key={schedule.id}>
                        {scheduleFields.map((field) => (
                            <ChakraTable.Cell
                                key={`${schedule.id}-${field.key}`}>
                                {format(schedule, field)}
                            </ChakraTable.Cell>
                        ))}
                    </ChakraTable.Row>
                ))}
            </ChakraTable.Body>
        </ChakraTable.Root>
    );
}

function format(schedule: Schedule, field: ScheduleField) {
    const value = schedule[field.key]?.toString() || '';

    if (field.format) {
        return field.format(value);
    }

    return value;
}

const scheduleFields: ScheduleField[] = [
    {key: 'title', value: "Title"},
    {key: 'description', value: "Description"},
    {key: 'status', value: "Status", format: statusFormat},
    {key: 'createdAt', value: "Created at", format: dateFormat},
    {key: 'runAt', value: "Run at", format: dateFormat},
    {key: 'startAt', value: "Start at", format: dateFormat},
    {key: 'endAt', value: "End at", format: dateFormat},
];

type ScheduleField = {
    key: keyof Schedule;
    value: string;
    format?: (v: string) => ReactNode | string;
};

function dateFormat(value: string): string {
    return dayjs(value).format('DD/MM/YYYY HH:mm:ss');
}

const SCHEDULE_STATUSES = [
    'not_started',
    'active',
    'paused',
    'expired'
]

function isScheduleStatus(value: string): value is ScheduleStatus {
    return SCHEDULE_STATUSES.includes(value);
}

const statusBadgeColorMap: Record<ScheduleStatus, ColorPalette> = {
    not_started: "gray",
    active: "green",
    paused: "yellow",
    expired: "red"
}

function statusFormat(status: string): ReactNode {
    if (!isScheduleStatus(status)) {
        return status;
    }

    return (
        <Badge textTransform="capitalize" colorPalette={statusBadgeColorMap[status]}>
            {status}
        </Badge>
    )
}