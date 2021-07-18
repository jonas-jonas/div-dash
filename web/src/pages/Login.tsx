import ky from "ky";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useHistory } from "react-router";

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
  const [error, setError] = useState<string | null>();
  const history = useHistory();

  const onSubmit = async (values: LoginForm) => {
    try {
      await ky.post("/api/login", {
        json: values,
      });

      history.replace("/");
    } catch (error) {
      if (error instanceof ky.HTTPError) {
        const errResponse: ApiError = await error.response.json();
        setError(errResponse.message);
        reset({ email: values.email, password: "" });
      }
    }
  };

  const dismissError = () => {
    setError(null);
  };

  return (
    <div className="container mx-auto">
      <form
        className="bg-white rounded shadow w-80 px-6 py-12 mx-auto mt-20"
        onSubmit={handleSubmit(onSubmit)}
      >
        <h1 className="text-3xl mb-8">Login</h1>
        <label className="block mb-4">
          <span className="text-xs text-gray-600 ml-3">Email</span>
          <input
            type="text"
            className="bg-gray-100 block w-full px-3 py-2 focus:outline-none rounded-md border border-gray-400 focus:border-blue-700 transition-colors shadow-inner"
            {...register("email", { required: true })}
          />
        </label>
        <label className="block mb-6">
          <span className="text-xs text-gray-600 ml-3">Password</span>
          <input
            type="password"
            className="bg-gray-100 block w-full px-3 py-2 focus:outline-none rounded-md border border-gray-400 focus:border-blue-700 transition-colors shadow-inner"
            {...register("password", { required: true })}
          />
        </label>
        {error && (
          <div className="bg-red-300 rounded-md py-2 px-4 border border-red-400 mb-4 flex justify-between">
            {error}
            <button className="text-red-700" onClick={dismissError}>
              x
            </button>
          </div>
        )}
        <button
          className="mx-auto block bg-gray-900 text-white rounded px-6 py-2 shadow hover:bg-gray-600 transition-colors focus:outline-none"
          type="submit"
        >
          Login
        </button>
      </form>
    </div>
  );
}
