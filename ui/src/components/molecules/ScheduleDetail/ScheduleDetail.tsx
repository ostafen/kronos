import Schedule from '@/model/schedule.ts';
import ScheduleStatusBadge from '@/components/atoms/ScheduleStatusBadge/ScheduleStatusBadge.tsx';
import formatDate from '@/utils/format-date.ts';
import { DataListItem, DataListRoot } from '@/components/chakra/data-list.tsx';
import { Flex } from '@chakra-ui/react';

function ScheduleDetail(schedule: Schedule) {
  const {
    id,
    title,
    description,
    metadata,
    isRecurring,
    cronExpr,
    url,
    runAt,
    createdAt,
    startAt,
    endAt,
    status,
  } = schedule;

  const items = [
    {
      label: 'ID',
      value: id,
    },
    {
      label: 'Status',
      value: <ScheduleStatusBadge status={status} />,
    },
    {
      label: 'Title',
      value: title,
    },
    {
      label: 'Description',
      value: description,
    },
    {
      label: 'Webhook URL',
      value: url,
    },
    {
      label: 'Created at',
      value: formatDate(createdAt),
    },
    ...(isRecurring
      ? [
          {
            label: 'Start at',
            value: formatDate(startAt),
          },
          {
            label: 'End at',
            value: formatDate(endAt),
          },
          {
            label: 'Cron expression',
            value: <code>{cronExpr}</code>,
          },
        ]
      : [
          {
            label: 'Run at',
            value: formatDate(runAt),
          },
        ]),
    {
      label: 'Metadata',
      value: metadata,
    },
  ];

  return (
    <>
      <DataListRoot orientation="horizontal" divideY="1px" maxW="md">
        {items.map((item) => (
          <DataListItem
            pt="4"
            grow
            wordBreak="break-word"
            key={item.label}
            label={item.label}
            value={item.value || '–'}
          />
        ))}
      </DataListRoot>

      <Flex></Flex>
    </>
  );
}

export default ScheduleDetail;
