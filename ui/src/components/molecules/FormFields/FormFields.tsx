import FormField from '@/model/form-field';
import { useWatch } from 'react-hook-form';
import InputFormField from '../InputFormField/InputFormField';
import { useMemo } from 'react';

export default function FormFields(props: { fields: FormField[] }) {
  const { fields } = props;

  const isRecurring = useWatch({
    name: 'isRecurring',
  });

  const filteredFields = useMemo(
    () =>
      fields
        .sort((a, b) => a.order - b.order)
        .filter((field) => {
          switch (field.when) {
            case 'not-recurring':
              return !isRecurring;

            case 'recurring':
              return isRecurring;

            default:
              return true;
          }
        }),
    [fields, isRecurring]
  );

  return (
    <>
      {filteredFields.map((field) => (
        <InputFormField {...field} key={field.order} />
      ))}
    </>
  );
}
