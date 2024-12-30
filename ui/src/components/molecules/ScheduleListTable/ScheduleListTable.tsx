import Schedule from '@/model/schedule';
import { Flex, Table as ChakraTable } from '@chakra-ui/react';
import dayjs from 'dayjs';
import { ReactNode } from 'react';
import { FaRegArrowAltCircleRight } from 'react-icons/fa';
import IconButtonLink from '@/components/atoms/IconButtonLink/IconButtonLink.tsx';
import ScheduleStatusBadge from '@/components/atoms/ScheduleStatusBadge/ScheduleStatusBadge.tsx';
import DeleteScheduleTrigger from '@/components/molecules/DeleteScheduleTrigger/DeleteScheduleTrigger.tsx';

interface SchedulesTableProps {
  schedules: Schedule[];
}

export default function SchedulesTable(props: SchedulesTableProps) {
  const { schedules } = props;

  return (
    <>
      <ChakraTable.Root>
        <ChakraTable.Header>
          <ChakraTable.Row>
            {scheduleFields.map((field) => (
              <ChakraTable.ColumnHeader key={field.key}>
                {field.value}
              </ChakraTable.ColumnHeader>
            ))}
            <ChakraTable.ColumnHeader>Actions</ChakraTable.ColumnHeader>
          </ChakraTable.Row>
        </ChakraTable.Header>
        <ChakraTable.Body>
          {schedules.map((schedule) => (
            <ChakraTable.Row key={schedule.id}>
              {scheduleFields.map((field) => (
                <ChakraTable.Cell key={`${schedule.id}-${field.key}`}>
                  {format(schedule, field)}
                </ChakraTable.Cell>
              ))}
              <ChakraTable.Cell>
                <Flex>
                  <DeleteScheduleTrigger scheduleId={schedule.id} />
                  <IconButtonLink
                    to={`/schedule/${schedule.id}`}
                    title="View schedule detail"
                    aria-label="View schedule detail"
                    variant="ghost"
                  >
                    <FaRegArrowAltCircleRight />
                  </IconButtonLink>
                </Flex>
              </ChakraTable.Cell>
            </ChakraTable.Row>
          ))}
        </ChakraTable.Body>
      </ChakraTable.Root>
    </>
  );
}

function format(schedule: Schedule, field: ScheduleField) {
  const value = schedule[field.key]?.toString() || '';

  if (field.format) {
    return field.format(value);
  }

  return value;
}

const scheduleFields: ScheduleField[] = [
  { key: 'title', value: 'Title' },
  { key: 'description', value: 'Description' },
  { key: 'status', value: 'Status', format: statusFormat },
  { key: 'createdAt', value: 'Created at', format: dateFormat },
  { key: 'runAt', value: 'Run at', format: dateFormat },
  { key: 'startAt', value: 'Start at', format: dateFormat },
  { key: 'endAt', value: 'End at', format: dateFormat },
];

type ScheduleField = {
  key: keyof Schedule;
  value: string;
  format?: (v: string) => ReactNode | string;
};

function dateFormat(value: string): string {
  return dayjs(value).format('DD/MM/YYYY HH:mm:ss');
}

function statusFormat(status: string) {
  return <ScheduleStatusBadge status={status} />;
}
