import { LuPause, LuPlay } from 'react-icons/lu';
import { Button, ButtonProps } from '@/components/chakra/button.tsx';
import useInvalidateSchedules from '@/hooks/use-invalidate-schedules.ts';
import useResumeSchedule from '@/hooks/use-resume-schedule.ts';
import usePauseSchedule from '@/hooks/use-pause-schedule.ts';
import { UseMutateAsyncFunction } from '@tanstack/react-query';
import useTriggerSchedule from '@/hooks/use-trigger-schedule.ts';
import { GrTrigger } from 'react-icons/gr';
import { ReactNode } from 'react';

export type ScheduleAction = 'pause' | 'resume' | 'trigger';

interface ActionConfig {
  handler: UseMutateAsyncFunction<Response, Error, string, unknown>;
  content: ReactNode;
}

interface ScheduleActionButtonProps extends ButtonProps {
  scheduleId: string[] | string;
  action: ScheduleAction;
}

function ScheduleActionButton(props: ScheduleActionButtonProps) {
  const { scheduleId, action, ...rest } = props;

  const resumeSchedule = useResumeSchedule();
  const pauseSchedule = usePauseSchedule();
  const triggerSchedule = useTriggerSchedule();

  const config: Record<ScheduleAction, ActionConfig> = {
    resume: {
      handler: resumeSchedule.mutateAsync,
      content: (
        <>
          <LuPlay />
          Resume
        </>
      ),
    },
    pause: {
      handler: pauseSchedule.mutateAsync,
      content: (
        <>
          <LuPause />
          Pause
        </>
      ),
    },
    trigger: {
      handler: triggerSchedule.mutateAsync,
      content: (
        <>
          <GrTrigger />
          Trigger
        </>
      ),
    },
  };

  const handleClick = useInvalidateSchedules(() => {
    const schedules = scheduleId instanceof Array ? scheduleId : [scheduleId];

    return Promise.allSettled(
      schedules.map((schedule) => config[action].handler(schedule))
    );
  });

  return (
    <Button {...rest} onClick={handleClick} variant="outline" size="sm">
      {config[action].content}
    </Button>
  );
}

export default ScheduleActionButton;
