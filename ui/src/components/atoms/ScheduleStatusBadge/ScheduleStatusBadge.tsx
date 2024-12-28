import {ScheduleStatus} from "@/model/schedule.ts";
import {Badge, BadgeProps, ColorPalette} from "@chakra-ui/react";
import {ReactNode} from "react";

export default function ScheduleStatusBadge({status, ...badgeProps}: ScheduleStatusBadgeProps): ReactNode {
    if (!isScheduleStatus(status)) {
        return status;
    }

    return (
        <Badge
            textTransform="capitalize"
            colorPalette={statusBadgeColorMap[status]}
            {...badgeProps}
        >
            {status}
        </Badge>
    );
}

const statusBadgeColorMap: Record<ScheduleStatus, ColorPalette> = {
    not_started: "gray",
    active: "green",
    paused: "yellow",
    expired: "red",
};

interface ScheduleStatusBadgeProps extends BadgeProps {
    status: string;
}

const SCHEDULE_STATUSES = ["not_started", "active", "paused", "expired"];

function isScheduleStatus(value: string): value is ScheduleStatus {
    return SCHEDULE_STATUSES.includes(value);
}