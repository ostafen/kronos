import ButtonLink from "@/components/atoms/ButtonLink/ButtonLink";
import useFetchSchedules from "@/hooks/use-fetch-schedules";
import {Button, Heading, Text} from "@chakra-ui/react";
import {FiPlus} from "react-icons/fi";
import ScheduleCardList from "@/components/molecules/ScheduleCardList/ScheduleCardList.tsx";
import {TbSeeding} from "react-icons/tb";
import {useQueryClient} from "@tanstack/react-query";
import seedDatabase from "@/seed/seed.js";

export default function ScheduleListPage() {
    const schedules = useFetchSchedules();
    const queryClient = useQueryClient();
    const noSchedules = !schedules.isFetching && !schedules.data?.length;
    const handleSeed = async () => {
        await seedDatabase();
        await queryClient.invalidateQueries({queryKey: ['schedules']});
    }

    return (
        <>
            <Heading as="h1" fontSize="2rem" mb="3">
                Welcome to Kronos
            </Heading>
            <Text fontSize="18px" mb="6">Your cronjob dashboard. ⏳</Text>

            {noSchedules && (
                <Text id="new-schedule-description" mb="3">👇 Click this button to add your first cron job
                    schedule.</Text>
            )}

            <ButtonLink aria-describedby="new-schedule-description" to="/new-schedule" mb="12" ml="auto">
                <FiPlus aria-hidden="true"/>
                New schedule
            </ButtonLink>
            {noSchedules && (
                <>
                    <Text mb="3">👇 Or click this one to quickly seed the database and test the application.</Text>
                    <Button mb="12" variant="subtle" onClick={handleSeed}>
                        <TbSeeding/>
                        Seed with initial data
                    </Button>
                    <Text fontWeight="600">New schedules will be displayed here. </Text>
                </>
            )}
            <ScheduleCardList schedules={schedules.data || []}/>
        </>
    );
}
