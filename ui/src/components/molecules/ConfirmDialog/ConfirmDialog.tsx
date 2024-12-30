import { Button } from '@/components/chakra/button';
import {
  DialogBody,
  DialogCloseTrigger,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogRoot,
  DialogTitle,
} from '@/components/chakra/dialog';
import { DialogActionTrigger } from '@chakra-ui/react';
import { useEffect, useState } from 'react';
import {
  closeDialog$,
  dialogConfirm$,
  dialogOpen$,
  dialogReset$,
  isDialogOpen$,
} from '@/utils/modal.ts';
import DialogData from '@/model/dialog-data.ts';

export default function ConfirmDialog() {
  const [isOpen, setIsOpen] = useState(false);
  const [dialogData, setDialogData] = useState<DialogData | null>(null);

  const handleConfirm = () => {
    setDialogData((data) => {
      if (!data) throw new Error('Data is not defined');
      return { ...data, isConfirmed: true };
    });

    dialogConfirm$.sink.next();
  };

  useEffect(() => {
    isDialogOpen$.sink.next(isOpen);
  }, [isOpen]);

  const handleCancel = () => {
    if (dialogData?.isConfirmed) return;
    setIsOpen(false);
  };

  useEffect(() => {
    let timeoutId: number | null = null;

    const subs = [
      dialogReset$.source$.subscribe(() => setDialogData(null)),
      closeDialog$.source$.subscribe(() => {
        setIsOpen(false);
        timeoutId = window.setTimeout(() => {
          setDialogData(null);
          timeoutId = null;
        }, 1000);
      }),
      dialogOpen$.source$.subscribe((dialogData) => {
        if (timeoutId) {
          window.clearTimeout(timeoutId);
        }

        setDialogData(dialogData);
        setIsOpen(true);
      }),
    ];

    return () => subs.forEach((sub) => sub.unsubscribe());
  }, []);

  return (
    <DialogRoot open={isOpen} onOpenChange={handleCancel} placement="center">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{dialogData?.title}</DialogTitle>
        </DialogHeader>
        <DialogBody>{dialogData?.content}</DialogBody>
        <DialogFooter>
          <DialogActionTrigger asChild>
            <Button variant="outline">Cancel</Button>
          </DialogActionTrigger>
          <Button onClick={handleConfirm}>Confirm</Button>
        </DialogFooter>
        <DialogCloseTrigger />
      </DialogContent>
    </DialogRoot>
  );
}
