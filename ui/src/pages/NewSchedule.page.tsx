import { Container, Heading, Stack } from "@chakra-ui/react";
import ChakraBreadcrumbLink from "@/components/atoms/ChakraBreadcrumbLink/ChakraBreadcrumbLink";
import {
  BreadcrumbCurrentLink,
  BreadcrumbRoot,
} from "@/components/chakra/breadcrumb";
import { Button } from "@/components/chakra/button";
import useCreateSchedule from "@/hooks/use-create-schedule";
import FormField from "@/model/form-field";
import { NewSchedule } from "@/model/schedule";
import { FormProvider, SubmitHandler, useForm } from "react-hook-form";
import { FiChevronRight } from "react-icons/fi";
import { useNavigate } from "react-router";
import formFields from "../assets/new-schedule.json";
import FormFields from "@/components/molecules/FormFields/FormFields";

export default function NewSchedulePage() {
  const form = useForm<NewSchedule>({
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

  const fields = formFields as FormField[];
  const createSchedule = useCreateSchedule();
  const navigate = useNavigate();

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
              startAt: `${runAt}:00+01:00`,
              endAt: `${runAt}:00+01:00`,
              ...(runAt && { runAt: `${runAt}:00+01:00` }),
            }),
      });

      navigate("/");
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <Container>
      <BreadcrumbRoot separator={<FiChevronRight />} variant="underline">
        <ChakraBreadcrumbLink to="/">Home</ChakraBreadcrumbLink>
        <BreadcrumbCurrentLink>New Schedule</BreadcrumbCurrentLink>
      </BreadcrumbRoot>

      <Heading as="h1" mt="3" mb="6">
        New Schedule
      </Heading>

      <FormProvider {...form}>
        <form onSubmit={form.handleSubmit(submitHandler)}>
          <Stack gap="4" align="flex-start">
            <FormFields fields={fields} />
            <Button type="submit">Submit</Button>
          </Stack>
        </form>
      </FormProvider>
    </Container>
  );
}
