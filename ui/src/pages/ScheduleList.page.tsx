import ButtonLink from "@/components/atoms/ButtonLink/ButtonLink";
import SchedulesTable from "@/components/molecules/ScheduleListTable/ScheduleListTable";
import useFetchSchedules from "@/hooks/use-fetch-schedules";
import {Container, Heading} from "@chakra-ui/react";
import {FiPlus} from "react-icons/fi";

export default function ScheduleListPage() {
    const schedules = useFetchSchedules();

    return (
        <Container mt="10">
            <Heading as="h1" fontSize="2rem" mb="6">
                Schedule List
            </Heading>
            <ButtonLink to="/new-schedule" mb="5" ml="auto">
                <FiPlus/>
                New schedule
            </ButtonLink>
            <SchedulesTable schedules={schedules.data || []}/>
        </Container>
    );
}
