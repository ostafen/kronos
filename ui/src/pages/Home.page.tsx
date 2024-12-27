import ButtonLink from "@/components/atoms/ButtonLink/ButtonLink";
import SchedulesTable from "@/components/molecules/ScheduleListTable/ScheduleListTable";
import useFetchSchedules from "@/hooks/use-fetch-schedules";
import { Container, Heading } from "@chakra-ui/react";
import { FiPlus } from "react-icons/fi";

export default function HomePage() {
  const schedules = useFetchSchedules();

  return (
    <Container mt="10">
      <Heading as="h1" fontSize="36px" mb="12">
        Schedule List
      </Heading>
      <ButtonLink to="/new-schedule" mb="5" ml="auto">
        <FiPlus />
        New schedule
      </ButtonLink>
      <SchedulesTable schedules={schedules.data || []} />
    </Container>
  );
}
