import { Checkbox, CheckboxProps } from '@/components/chakra/checkbox.tsx';

interface ToggleAllCheckboxProps extends CheckboxProps {
  /**
   * The array of the checked items.
   */
  checkedItems: unknown[];

  /**
   * The array that contains all the check-able items.
   */
  allItems: unknown[];

  /**
   * Callback triggered when the checkbox is pressed when all the items are checked.
   */
  onResetCheckedItems(): void;

  /**
   * Callback triggered when the checkbox is pressed and not all the items are checked.
   */
  onCheckAllItems(): void;
}

export default function ToggleAllCheckbox(props: ToggleAllCheckboxProps) {
  const {
    checkedItems,
    allItems,
    onCheckAllItems,
    onResetCheckedItems,
    ...rest
  } = props;
  const areAllChecked = checkedItems.length === allItems.length;
  const indeterminate = checkedItems.length > 0 && !areAllChecked;

  const handleCheck = () => {
    if (areAllChecked) {
      return onResetCheckedItems();
    }

    onCheckAllItems();
  };

  return (
    <Checkbox
      onCheckedChange={handleCheck}
      checked={indeterminate ? 'indeterminate' : areAllChecked}
      {...rest}
    >
      Select all schedules
    </Checkbox>
  );
}
