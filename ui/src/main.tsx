import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Provider } from './components/chakra/provider.tsx';
import { BrowserRouter, Routes } from 'react-router';
import { Route } from 'react-router';
import NewSchedulePage from './pages/NewSchedule.page.tsx';
import HomePage from './pages/Home.page.tsx';
import RootLayout from './components/layouts/Root.layout.tsx';

const queryClient = new QueryClient();

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <Provider defaultTheme="light">
        <BrowserRouter basename="/web">
          <Routes>
            <Route element={<RootLayout />}>
              <Route index element={<HomePage />} />
              <Route path="/new-schedule" element={<NewSchedulePage />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </Provider>
    </QueryClientProvider>
  </StrictMode>,
);
