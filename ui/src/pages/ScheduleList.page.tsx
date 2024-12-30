import ButtonLink from "@/components/atoms/ButtonLink/ButtonLink";
import useFetchSchedules from "@/hooks/use-fetch-schedules";
import {Box, Button, Flex, Heading, Text} from "@chakra-ui/react";
import {FiPlus} from "react-icons/fi";
import ScheduleCardList from "@/components/molecules/ScheduleCardList/ScheduleCardList.tsx";
import {TbSeeding} from "react-icons/tb";
import {useQueryClient} from "@tanstack/react-query";
import seedDatabase from "@/seed/seed.js";
import {keyframes} from "@emotion/react";
import GradientText from "@/components/atoms/GradientText/GradientText.tsx";
import FadeInBox from "@/components/atoms/FadeInBox/FadeInBox.tsx";
import {useContext, useEffect} from "react";
import AppContext from "@/context/App.context.tsx";


const shinyBar = keyframes`
    from {
        left: 0;
    }

    to {
        left: calc(100%);
    }
`

export default function ScheduleListPage() {
    const schedules = useFetchSchedules();
    const queryClient = useQueryClient();

    const {state, dispatch} = useContext(AppContext);
    const noSchedules = !schedules.isFetching && !schedules.data?.length;
    const isAnimationDisabled = state.animationStatus === 'played';

    const handleSeed = async () => {
        await seedDatabase();
        await queryClient.invalidateQueries({queryKey: ['schedules']});
    }

    useEffect(() => {
        if (window.sessionStorage.getItem("animationStatus")) {
            dispatch({animationStatus: "played"});
        }

        window.sessionStorage.setItem("animationStatus", "played");
    }, []);

    return (
        <FadeInBox
            as={Flex}
            isAnimationDisabled={isAnimationDisabled}
            align="center"
            textAlign="center"
            flexDirection="column">
            <Heading as="h1"
                     overflow="hidden"
                     mb="4" mt="32"
                     fontSize="4rem"
                     lineHeight="4rem"
                     position="relative"
                     {...noSchedules && {
                         fontSize: {base: "4rem", lg: "6.4rem"},
                         lineHeight: {base: "4rem", lg: "6.4rem"},
                         height: "auto",
                     }}>
                <Box
                    left={0}
                    top="-1rem"
                    transform="rotate(12deg)"
                    position="absolute"
                    height="8rem"
                    width="4rem"
                    bgImage="linear-gradient(to right, transparent, white 50%, transparent)"
                    {...!isAnimationDisabled && {animation: `${shinyBar} 3s ease-in-out`}}
                />
                <GradientText fontWeight="900" text="Welcome to Kronos"/>
            </Heading>

            <FadeInBox isAnimationDisabled={isAnimationDisabled}
                       as={Text}
                       animationDelay="1s"
                       fontSize="1.8rem"
                       lineHeight="3.6rem"
                       mb="10">
                <GradientText as="strong" fontWeight="700" text="Your"/>{" "}
                cron dashboard
            </FadeInBox>

            <FadeInBox isAnimationDisabled={isAnimationDisabled}
                       animationDelay="3s">
                {noSchedules && <Heading mb="5">Get started, right now.</Heading>}

                <Box mb="8">
                    <Text mb="4" id="new-schedule-description">
                        {noSchedules ?
                            "Click the button below to add your first cron job schedule." :
                            "Add a new cron job schedule ⏳👇"}
                    </Text>

                    <ButtonLink display="flex"
                                justifySelf="center"
                                aria-describedby="new-schedule-description"
                                to="/new-schedule">
                        <FiPlus aria-hidden="true"/>
                        New schedule
                    </ButtonLink>
                </Box>

                {noSchedules && <Box mb="16">
                    <Text mb="4">Or click this one to quickly seed the database and test the application.</Text>
                    <Button variant="subtle" onClick={handleSeed}>
                        <TbSeeding/>
                        Seed with initial data
                    </Button>
                </Box>}
            </FadeInBox>

            {noSchedules && (
                <FadeInBox isAnimationDisabled={isAnimationDisabled}
                           animationDelay="4s"
                           as={Text}
                           fontSize="1.2rem"
                           fontWeight="500">New schedules
                    will be displayed here.</FadeInBox>
            )}

            {!noSchedules && (
                <FadeInBox isAnimationDisabled={isAnimationDisabled}>
                    <Box as="hr" mb="16" boxShadow="0 0 1px 0 rgba(212,212,212,.8)" width="100%"/>
                    <ScheduleCardList schedules={schedules.data || []}/>
                </FadeInBox>)}
        </FadeInBox>
    );
}
