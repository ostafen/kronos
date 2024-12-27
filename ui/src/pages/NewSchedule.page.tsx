import { Field } from "@/components/chakra/field";
import { Container, Heading, Input, Stack, Textarea } from "@chakra-ui/react";

import ChakraBreadcrumbLink from "@/components/atoms/ChakraBreadcrumbLink/ChakraBreadcrumbLink";
import {
  BreadcrumbCurrentLink,
  BreadcrumbRoot,
} from "@/components/chakra/breadcrumb";
import { Button } from "@/components/chakra/button";
import { Controller, useForm, useWatch } from "react-hook-form";
import { FiChevronRight } from "react-icons/fi";
import { Switch } from "@/components/ui/switch";
import useCreateSchedule from "@/hooks/use-create-schedule";
import { useNavigate } from "react-router";
import { NewSchedule } from "@/model/schedule";

export default function NewSchedulePage() {
  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
  } = useForm<NewSchedule>({
    defaultValues: {
      title: "",
      description: "",
      startAt: "",
      runAt: "",
      endAt: "",
      cronExpr: "",
      url: "",
      isRecurring: false,
    },
  });

  const isRecurring = useWatch({
    control,
    name: "isRecurring",
  });

  const createSchedule = useCreateSchedule();
  const navigate = useNavigate();

  const onSubmit = handleSubmit(async (data) => {
    const { startAt, endAt, runAt, isRecurring, ...otherData } = data;

    await createSchedule.mutateAsync({
      ...otherData,
      isRecurring: isRecurring || false,
      ...(isRecurring
        ? {
            ...(startAt && { startAt: `${startAt}:00+01:00` }),
            ...(endAt && { endAt: `${endAt}:00+01:00` }),
          }
        : {
            startAt: `${runAt}:00+01:00`,
            endAt: `${runAt}:00+01:00`,
            ...(runAt && { runAt: `${runAt}:00+01:00` }),
          }),
    });

    navigate("/");
  });

  return (
    <Container>
      <BreadcrumbRoot separator={<FiChevronRight />} variant="underline">
        <ChakraBreadcrumbLink to="/">Home</ChakraBreadcrumbLink>
        <BreadcrumbCurrentLink>New Schedule</BreadcrumbCurrentLink>
      </BreadcrumbRoot>
      <Heading mt="3" mb="6">
        New Schedule
      </Heading>
      <form onSubmit={onSubmit}>
        <Stack gap="4" align="flex-start">
          <Field
            label="Title"
            required
            invalid={!!errors.title}
            errorText={errors.title?.message}
          >
            <Input {...register("title", { required: "Title is required" })} />
          </Field>
          <Field
            label="Description"
            required
            invalid={!!errors.description}
            errorText={errors.description?.message}
          >
            <Textarea
              {...register("description", {
                required: "Description is required",
              })}
            />
          </Field>
          <Field
            label="Webhook URL"
            required
            invalid={!!errors.url}
            errorText={errors.url?.message}
          >
            <Input
              {...register("url", {
                required: "Webhook URL is required",
              })}
            />
          </Field>
          {!isRecurring ? (
            <Field
              label="Run at"
              required
              invalid={!!errors.runAt}
              errorText={errors.runAt?.message}
            >
              <Input
                type="datetime-local"
                {...register("runAt", {
                  required: "Run at is required",
                })}
              />
            </Field>
          ) : (
            <>
              <Field
                label="Start at"
                required
                invalid={!!errors.startAt}
                errorText={errors.startAt?.message}
              >
                <Input
                  type="datetime-local"
                  {...register("startAt", {
                    required: "Start at is required",
                  })}
                />
              </Field>
              <Field
                label="End at"
                required
                invalid={!!errors.url}
                errorText={errors.url?.message}
              >
                <Input
                  type="datetime-local"
                  {...register("endAt", {
                    required: "End at is required",
                  })}
                />
              </Field>
            </>
          )}
          <Field
            label="Cron expression"
            required
            invalid={!!errors.url}
            errorText={errors.url?.message}
          >
            <Input
              {...register("cronExpr", {
                required: "Cron expression is required",
              })}
            />
          </Field>

          <Controller
            name="isRecurring"
            control={control}
            render={({ field }) => (
              <Field
                invalid={!!errors.isRecurring}
                errorText={errors.isRecurring?.message}
              >
                <Switch
                  name={field.name}
                  checked={field.value}
                  onCheckedChange={({ checked }) => field.onChange(checked)}
                  inputProps={{ onBlur: field.onBlur }}
                >
                  Is recurring
                </Switch>
              </Field>
            )}
          />
          <Button type="submit">Submit</Button>
        </Stack>
      </form>
    </Container>
  );
}
