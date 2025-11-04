// imports libraries(components)
import { Outlet } from "react-router-dom";

// creates and exports layout component
const Layout = () => {
  return (
    <main>
      {/* Outlet - is used to render child routes inside the layout */}
      <Outlet />
    </main>
  );
};

export default Layout;
