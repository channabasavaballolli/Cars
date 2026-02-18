import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";
import { graphqlRequest, REQUEST_LOGIN, VERIFY_LOGIN } from "@/lib/graphql";

const AdminLogin = () => {
    const [email, setEmail] = useState("");
    const [otp, setOtp] = useState("");
    const [step, setStep] = useState<"email" | "otp">("email");
    const [loading, setLoading] = useState(false);
    const { login } = useAuth();
    const navigate = useNavigate();

    const handleRequestLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        try {
            await graphqlRequest(REQUEST_LOGIN, { email });
            setStep("otp");
            toast.success("Verification code sent to your email");
        } catch (error) {
            toast.error("Failed to send verification code");
        } finally {
            setLoading(false);
        }
    };

    const handleVerifyLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        try {
            const data = await graphqlRequest(VERIFY_LOGIN, { email, code: otp });
            login(data.verifyLogin);
            toast.success("Welcome back, Admin!");
            navigate("/dashboard");
        } catch (error) {
            toast.error("Invalid verification code");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white p-4">
            <div className="w-full max-w-md space-y-8 bg-gray-800 p-8 rounded-xl shadow-2xl border border-gray-700">
                <div className="text-center">
                    <h2 className="text-3xl font-bold bg-gradient-to-r from-red-500 to-orange-500 bg-clip-text text-transparent">
                        Admin Portal
                    </h2>
                    <p className="mt-2 text-gray-400">Secure access for fleet management</p>
                </div>

                {step === "email" ? (
                    <form onSubmit={handleRequestLogin} className="space-y-6">
                        <div className="space-y-2">
                            <label htmlFor="email" className="block text-sm font-medium text-gray-300">
                                Admin Email
                            </label>
                            <Input
                                id="email"
                                type="email"
                                placeholder="admin@company.com"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                required
                                className="bg-gray-900 border-gray-700 text-white placeholder-gray-500 focus:border-red-500 focus:ring-red-500"
                            />
                        </div>
                        <Button
                            type="submit"
                            className="w-full bg-gradient-to-r from-red-600 to-orange-600 hover:from-red-700 hover:to-orange-700 text-white font-bold py-2 px-4 rounded transition-all duration-200"
                            disabled={loading}
                        >
                            {loading ? "Sending Code..." : "Continue"}
                        </Button>
                        <div className="text-center mt-4">
                            <a href="/login" className="text-sm text-gray-400 hover:text-white transition-colors">Not an admin? Go to User Login</a>
                        </div>
                    </form>
                ) : (
                    <form onSubmit={handleVerifyLogin} className="space-y-6">
                        <div className="space-y-2">
                            <label htmlFor="otp" className="block text-sm font-medium text-gray-300">
                                Verification Code
                            </label>
                            <Input
                                id="otp"
                                type="text"
                                placeholder="123456"
                                value={otp}
                                onChange={(e) => setOtp(e.target.value)}
                                required
                                maxLength={6}
                                className="bg-gray-900 border-gray-700 text-white placeholder-gray-500 focus:border-red-500 focus:ring-red-500 tracking-widest text-center text-2xl"
                            />
                        </div>
                        <Button
                            type="submit"
                            className="w-full bg-gradient-to-r from-red-600 to-orange-600 hover:from-red-700 hover:to-orange-700 text-white font-bold py-2 px-4 rounded transition-all duration-200"
                            disabled={loading}
                        >
                            {loading ? "Verifying..." : "Access Dashboard"}
                        </Button>
                        <button
                            type="button"
                            onClick={() => setStep("email")}
                            className="w-full text-sm text-gray-400 hover:text-white transition-colors"
                        >
                            Back to Email
                        </button>
                    </form>
                )}
            </div>
        </div>
    );
};

export default AdminLogin;
