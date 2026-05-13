import { createContext, useContext, useState, useCallback, type ReactNode } from 'react';

interface SidebarExtraCtx {
  extra: ReactNode | null;
  extraTitle: string;
  setExtra: (content: ReactNode | null, title?: string) => void;
}

const Ctx = createContext<SidebarExtraCtx>({
  extra: null,
  extraTitle: 'Contents',
  setExtra: () => {},
});

export function useSidebarExtra() {
  return useContext(Ctx);
}

export function SidebarExtraProvider({ children }: { children: ReactNode }) {
  const [extra, setExtraState] = useState<ReactNode | null>(null);
  const [extraTitle, setExtraTitle] = useState('Contents');

  const setExtra = useCallback((content: ReactNode | null, title?: string) => {
    setExtraState(content);
    if (title !== undefined) setExtraTitle(title);
  }, []);

  return (
    <Ctx.Provider value={{ extra, extraTitle, setExtra }}>
      {children}
    </Ctx.Provider>
  );
}
