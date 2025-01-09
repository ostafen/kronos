import { FiPlus } from 'react-icons/fi';
import DialogActionTrigger from '@/components/molecules/DialogActionTrigger/DialogActionTrigger.tsx';
import useCreateSchedule from '@/hooks/use-create-schedule.ts';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { NewSchedule } from '@/model/schedule.ts';
import { Stack } from '@chakra-ui/react';
import FormFields from '@/components/molecules/FormFields/FormFields.tsx';
import formFields from '@/assets/new-schedule.json';
import FormField from '@/model/form-field.ts';
import { useEffect } from 'react';
import { isDialogOpen$ } from '@/utils/modal.ts';
import { useQueryClient } from '@tanstack/react-query';

export default function NewScheduleTrigger() {
  const form = useForm<NewSchedule>({
    defaultValues: {
      title: '',
      description: '',
      startAt: '',
      runAt: '',
      endAt: '',
      cronExpr: '',
      url: '',
      isRecurring: false,
    },
  });

  const queryClient = useQueryClient();
  const fields = formFields as FormField[];
  const createSchedule = useCreateSchedule();

  useEffect(() => {
    const sub = isDialogOpen$.source$.subscribe((isOpen) => {
      if (!isOpen) {
        setTimeout(() => form.reset(), 500);
      }
    });

    return () => sub.unsubscribe();
  }, []);

  const submitHandler: SubmitHandler<NewSchedule> = async (data) => {
    const { cronExpr, startAt, endAt, runAt, isRecurring, ...otherData } = data;

    try {
      await createSchedule.mutateAsync({
        ...otherData,
        isRecurring,
        ...(isRecurring
          ? {
              cronExpr,
              ...(startAt && { startAt: `${startAt}:00+01:00` }),
              ...(endAt && { endAt: `${endAt}:00+01:00` }),
            }
          : {
              ...(runAt && { runAt: `${runAt}:00+01:00` }),
            }),
      });
    } catch (error) {
      console.error(error);
    }
  };

  const dialogData = {
    title: 'Add new schedule',
    content: (
      <FormProvider {...form}>
        <Stack gap="4" align="flex-start">
          <FormFields fields={fields} />
        </Stack>
      </FormProvider>
    ),
  };

  return (
    <DialogActionTrigger
      onConfirm={form.handleSubmit(submitHandler)}
      onSuccess={() =>
        queryClient.invalidateQueries({ queryKey: ['schedules'] })
      }
      dialogData={dialogData}
    >
      <FiPlus aria-hidden="true" />
      New schedule
    </DialogActionTrigger>
  );
}
