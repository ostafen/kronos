export default interface FormField {
  label: string;
  isRequired: boolean;
  type: FormFieldType;
  name: string;
  when: ShowConditionType;
  order: number;
}

export type ShowConditionType = 'recurring' | 'not-recurring';

export type FormFieldType =
  | 'input-text'
  | 'input-url'
  | 'input-datetime-local'
  | 'input-cron'
  | 'textarea'
  | 'checkbox';
