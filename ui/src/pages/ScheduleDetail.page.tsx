import {Badge, Container, Flex, Heading, Text} from "@chakra-ui/react";
import {useNavigate, useParams} from "react-router";
import useFetchSchedule from "@/hooks/use-fetch-schedule.ts";
import {FiChevronRight} from "react-icons/fi";
import ChakraBreadcrumbLink from "@/components/atoms/ChakraBreadcrumbLink/ChakraBreadcrumbLink.tsx";
import {BreadcrumbCurrentLink, BreadcrumbRoot} from "@/components/chakra/breadcrumb.tsx";
import ScheduleStatusBadge from "@/components/atoms/ScheduleStatusBadge/ScheduleStatusBadge.tsx";
import IconButtonLink from "@/components/atoms/IconButtonLink/IconButtonLink.tsx";
import {RiExternalLinkLine} from "react-icons/ri";
import formatDate from "@/utils/format-date.ts";
import DeleteScheduleTrigger from "@/components/molecules/DeleteScheduleTrigger/DeleteScheduleTrigger.tsx";

export default function ScheduleDetailPage() {
    const {scheduleId} = useParams();
    const navigate = useNavigate();
    const schedule = useFetchSchedule(scheduleId);

    if (!schedule.data || schedule.error) {
        navigate('/');
        return null;
    }

    const {
        title,
        description,
        status,
        cronExpr,
        startAt,
        endAt,
        createdAt,
        runAt,
        metadata,
        isRecurring,
        id,
        url
    } = schedule.data;

    return (
        <Container mt="10">
            <BreadcrumbRoot separator={<FiChevronRight/>} variant="underline">
                <ChakraBreadcrumbLink to="/">Home</ChakraBreadcrumbLink>
                <BreadcrumbCurrentLink>Schedule Detail</BreadcrumbCurrentLink>
            </BreadcrumbRoot>

            <Flex direction="column" mt="10" gap="3">
                <Flex gap="3" align="center">
                    <Heading fontSize="2rem">{title}</Heading>
                    <IconButtonLink variant="plain" minW={0} to={url}>
                        <RiExternalLinkLine/>
                    </IconButtonLink>
                </Flex>
                <Flex gap="2" alignItems="center" flexWrap="wrap" maxW={{base: "100%", lg: "60%"}}>
                    <Badge title="Schedule unique id" colorPalette="orange">{id}</Badge>
                    <ScheduleStatusBadge alignSelf="flex-start" status={status}/>
                    {isRecurring && <Badge colorPalette="purple">Recurring</Badge>}
                    {cronExpr && <Badge colorPalette="red">{cronExpr}</Badge>}
                    <Badge colorPalette="cyan">Created at {formatDate(createdAt)}</Badge>
                    {runAt && new Date(runAt).getFullYear() > 1 &&
                        <Badge colorPalette="cyan">Run at {formatDate(runAt)}</Badge>}
                    {startAt && <Badge colorPalette="cyan">Start at {formatDate(startAt)}</Badge>}
                    {endAt && <Badge colorPalette="cyan">End at {formatDate(endAt)}</Badge>}
                </Flex>

                <Flex mt="1" mb="6" gap="4">
                    <DeleteScheduleTrigger
                        onSuccess={() => navigate('/')}
                        scheduleId={id}
                        colorPalette="red"
                        variant="plain"
                    >
                        Delete schedule
                    </DeleteScheduleTrigger>
                </Flex>

                <Text>{description}</Text>

                {metadata && <pre>{JSON.stringify(metadata, null, 2)}</pre>}
            </Flex>
        </Container>
    );
}
