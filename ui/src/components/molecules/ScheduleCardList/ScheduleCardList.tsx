import Schedule from '@/model/schedule.ts';
import { CheckboxCard } from '@/components/chakra/checkbox-card.tsx';
import {
  Badge,
  CheckboxGroup,
  Fieldset,
  Flex,
  Grid,
  GridItem,
  HStack,
} from '@chakra-ui/react';
import {
  ActionBarContent,
  ActionBarRoot,
  ActionBarSelectionTrigger,
  ActionBarSeparator,
} from '@/components/chakra/action-bar.tsx';
import { useEffect, useMemo, useState } from 'react';
import DeleteScheduleTrigger from '@/components/molecules/DeleteScheduleTrigger/DeleteScheduleTrigger.tsx';
import ScheduleStatusBadge from '@/components/atoms/ScheduleStatusBadge/ScheduleStatusBadge.tsx';
import DialogActionTrigger from '@/components/molecules/DialogActionTrigger/DialogActionTrigger.tsx';
import ScheduleActionButton, {
  ScheduleAction,
} from '@/components/atoms/ScheduleActionButton/ScheduleActionButton.tsx';
import ScheduleDetail from '@/components/molecules/ScheduleDetail/ScheduleDetail.tsx';

interface ScheduleCardListProps {
  schedules: Schedule[];
}

export default function ScheduleCardList(props: ScheduleCardListProps) {
  const { schedules } = props;
  const [checkedSchedulesIds, setCheckedSchedulesIds] = useState<string[]>([]);

  useEffect(() => {
    setTimeout(() => {
      setCheckedSchedulesIds([]);
    }, 200);
  }, [schedules]);

  const checkedSchedules: Schedule[] = useMemo(() => {
    const scheduleMap = schedules.reduce(
      (acc, schedule) => ({
        ...acc,
        [schedule.id]: schedule,
      }),
      {}
    );

    return checkedSchedulesIds
      .map((id) => scheduleMap[id as keyof typeof scheduleMap])
      .filter(Boolean);
  }, [schedules, checkedSchedulesIds]);

  const isPauseButtonDisabled = checkedSchedules.some(
    (schedule) => !schedule.isRecurring || schedule.status === 'paused'
  );

  const isResumeButtonDisabled = checkedSchedules.some(
    (schedule) => !schedule.isRecurring || schedule.status === 'active'
  );

  const scheduleActions: { action: ScheduleAction; isDisabled: boolean }[] = [
    {
      action: 'pause',
      isDisabled: isPauseButtonDisabled,
    },
    {
      action: 'resume',
      isDisabled: isResumeButtonDisabled,
    },
    {
      action: 'trigger',
      isDisabled: false,
    },
  ];

  return (
    <>
      {schedules.length > 0 && (
        <Fieldset.Root>
          <CheckboxGroup
            value={checkedSchedulesIds}
            onValueChange={(schedules) => setCheckedSchedulesIds(schedules)}
          >
            <Fieldset.Legend textAlign="left" fontSize="sm" mb="2">
              Schedules
            </Fieldset.Legend>
            <Grid
              templateColumns={{
                base: '1fr',
                md: '1fr 1fr',
                lg: 'repeat(3, 1fr)',
              }}
              gap="4"
            >
              {schedules.map((schedule) => (
                <GridItem key={schedule.id}>
                  <CheckboxCard
                    backgroundColor="white"
                    variant="outline"
                    colorPalette="purple"
                    h="100%"
                    label={schedule.title}
                    description={schedule.description}
                    value={schedule.id}
                    addon={
                      <Flex justify="space-between" align="center">
                        <DialogActionTrigger
                          variant="plain"
                          p={0}
                          onConfirm={() => Promise.resolve()}
                          dialogData={{
                            title: 'Details',
                            hideFooterButtons: true,
                            content: <ScheduleDetail {...schedule} />,
                          }}
                        >
                          Show details
                        </DialogActionTrigger>
                        <HStack>
                          {schedule.isRecurring && (
                            <>
                              <Badge colorPalette="purple">Recurring</Badge>
                              <Badge colorPalette="blue">
                                {schedule.cronExpr}
                              </Badge>
                            </>
                          )}
                          <ScheduleStatusBadge status={schedule.status} />
                        </HStack>
                      </Flex>
                    }
                  />
                </GridItem>
              ))}
            </Grid>
          </CheckboxGroup>
        </Fieldset.Root>
      )}

      <ActionBarRoot open={checkedSchedulesIds.length > 0}>
        <ActionBarContent>
          <ActionBarSelectionTrigger>
            {checkedSchedulesIds.length} selected
          </ActionBarSelectionTrigger>
          <ActionBarSeparator />
          <DeleteScheduleTrigger
            colorPalette="black"
            variant="outline"
            size="sm"
            scheduleId={checkedSchedulesIds}
          />
          {scheduleActions.map(({ action, isDisabled }) => (
            <ScheduleActionButton
              key={action}
              action={action}
              scheduleId={checkedSchedulesIds}
              disabled={isDisabled}
            />
          ))}
        </ActionBarContent>
      </ActionBarRoot>
    </>
  );
}
