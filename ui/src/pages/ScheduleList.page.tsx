import ButtonLink from "@/components/atoms/ButtonLink/ButtonLink";
import useFetchSchedules from "@/hooks/use-fetch-schedules";
import {Container, Heading} from "@chakra-ui/react";
import {FiPlus} from "react-icons/fi";
import ScheduleCardList from "@/components/molecules/ScheduleCardList/ScheduleCardList.tsx";

export default function ScheduleListPage() {
    const schedules = useFetchSchedules();

    return (
        <Container mt="10" mb="20">
            <Heading as="h1" fontSize="2rem" mb="6">
                Schedule List
            </Heading>
            <ButtonLink to="/new-schedule" mb="5" ml="auto">
                <FiPlus/>
                New schedule
            </ButtonLink>
            <ScheduleCardList schedules={schedules.data || []}/>
        </Container>
    );
}
