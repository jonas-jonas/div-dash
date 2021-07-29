import ky from "ky";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useHistory } from "react-router";
import { Link } from "react-router-dom";
import { useRecoilState } from "recoil";
import { User } from "../models/user";
import { userState } from "../state/authState";

type LoginForm = {
  email: string;
  password: string;
};

type ApiError = {
  message: string;
  status: number;
  path: string;
  timestamp: string;
};

export function Login() {
  const { register, handleSubmit, reset } = useForm<LoginForm>();
  const [, setUser] = useRecoilState(userState);
  const [error, setError] = useState<string | null>();
  const history = useHistory();

  const onSubmit = async (values: LoginForm) => {
    try {
      await ky.post("/api/login", {
        json: values,
      });
      
      const res = await ky.get("/api/auth/identity")
      const user: User = await res.json();
      setUser(user);

      history.replace("/");
    } catch (error) {
      if (error instanceof ky.HTTPError) {
        const errResponse: ApiError = await error.response.json();
        setError(errResponse.message);
        reset({ email: values.email, password: "" });
      }
    }
  };

  return (
    <div className="w-full h-screen relative">
      <div className="w-full sm:w-2/3 md:w-1/2 xl:w-1/3 bg-white h-full relative z-10">
        <form className="px-8 lg:px-20 py-8 mx-auto" onSubmit={handleSubmit(onSubmit)}>
          <div className="mb-24">
            <img src="/logo-dark@2x.png" alt="div-dash logo" />
          </div>
          <h1 className="text-3xl mb-3 font-semibold">Log In</h1>
          <h2 className="mb-14 lg:mb-20">Login with your email and password</h2>
          <label className="block mb-4">
            <span className="">Email</span>
            <input
              type="text"
              className="bg-gray-50 block w-full px-3 py-2 focus:outline-none rounded border border-transparent focus:border-blue-700 transition-colors"
              placeholder="Enter your email"
              {...register("email", { required: true })}
            />
          </label>
          <label className="block">
            <span className="">Password</span>
            <input
              type="password"
              className="bg-gray-50 block w-full px-3 py-2 focus:outline-none rounded border border-transparent focus:border-blue-700 transition-colors"
              placeholder="Enter your password"
              {...register("password", { required: true })}
            />
          </label>
          {error && (
            <div className="text-red-500 mt-2">
              {error}
            </div>
          )}
          <div className="flex justify-between items-center mt-6">
            <div>
              <Link to="/register" className="text-gray-500 mr-5">
                Register
              </Link>
              <Link to="/forgot-password" className="text-gray-500">
                Forgot Password?
              </Link>
            </div>
            <button
              className="bg-gray-900 text-white rounded px-6 py-2 shadow hover:bg-gray-600 transition-colors focus:outline-none"
              type="submit"
            >
              Login
            </button>
          </div>
        </form>
        <footer className="px-8 lg:px-20 py-8">
          <Link to="/impressum" className="mr-10 text-sm">
            Impressum
          </Link>
          <Link to="/contact" className="mr-10 text-sm">
            Contact Us
          </Link>
          <Link to="/about" className="mr-10 text-sm">
            About
          </Link>
        </footer>
      </div>
      <div className="fixed top-0 right-0 z-0 h-screen">
        <img className="bg-cover h-full" src="/login-bg.jpg" alt="bg" />
      </div>
    </div>
  );
}
