import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Provider } from './components/chakra/provider.tsx';
import ConfirmDialog from '@/components/molecules/ConfirmDialog/ConfirmDialog.tsx';
import FadeInBox from '@/components/atoms/FadeInBox/FadeInBox.tsx';
import { AppContextProvider } from '@/context/App.context.tsx';
import { Container } from '@chakra-ui/react';
import ScheduleList from '@/components/organisms/ScheduleList.tsx';

const queryClient = new QueryClient();

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <Provider defaultTheme="light">
        <AppContextProvider animationStatus="normal">
          {(state) => (
            <>
              <FadeInBox
                isAnimationDisabled={state.animationStatus === 'played'}
                zIndex={-1}
                width="100%"
                position="fixed"
                animationDelay="2s"
                minHeight="100vh"
                backgroundImage="radial-gradient(circle at 100%, transparent, transparent 50%, #efefef 75%, transparent 75%)"
              />

              <Container fluid pb="10">
                <ScheduleList />
              </Container>

              <ConfirmDialog />
            </>
          )}
        </AppContextProvider>
      </Provider>
    </QueryClientProvider>
  </StrictMode>
);
