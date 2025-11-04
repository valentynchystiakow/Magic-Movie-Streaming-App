// imports hooks
import { createContext, useState, useEffect } from "react";

// defines auth context
const AuthContext = createContext({});

// creates and exports auth provider component that manages authentication
export const AuthProvider = ({ children }) => {
  // hooks section
  const [auth, setAuth] = useState();
  const [loading, setLoading] = useState(true);

  // uses useEffect hook to fetch user data from local storage when the component mounts
  useEffect(() => {
    // try catch block to handle errors while parsing user data
    try {
      const storedUser = localStorage.getItem("user");
      if (storedUser) {
        const parsedUser = JSON.parse(storedUser);
        setAuth(parsedUser);
      }
    } catch (error) {
      console.error("Failed to parse user from localStorage", error);
    } finally {
      setLoading(false);
    }
  }, []);
  // uses useEffect hook to update user data in local storage when auth state changes
  useEffect(() => {
    if (auth) {
      localStorage.setItem("user", JSON.stringify(auth));
    } else {
      localStorage.removeItem("user");
    }
  }, [auth]);

  return (
    // wraps children with auth context provider to provide access to auth state
    <AuthContext.Provider value={{ auth, setAuth, loading }}>
      {children}
    </AuthContext.Provider>
  );
};
export default AuthContext;
