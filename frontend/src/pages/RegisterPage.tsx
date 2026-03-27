import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { authService } from '../services/authService';
import { useAuthStore, type AuthState } from '../store/authStore';
import toast from 'react-hot-toast';
import { Mail, Lock, User, ArrowRight, ShoppingBag, Gift, Heart, Star } from 'lucide-react';
import Button from '../components/ui/Button';
import Input from '../components/ui/Input';

export default function RegisterPage() {
  const navigate = useNavigate();
  const setAuth = useAuthStore((s: AuthState) => s.setAuth);
  const [form, setForm] = useState({ email: '', password: '', first_name: '', last_name: '' });
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const { user, tokens } = await authService.register(form);
      setAuth(user, tokens.access_token);
      toast.success('Account created!');
      navigate('/');
    } catch {
      toast.error('Registration failed');
    } finally {
      setLoading(false);
    }
  };

  const update = (field: string) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm({ ...form, [field]: e.target.value });

  const passwordStrength = form.password.length === 0 ? 0 : form.password.length < 6 ? 1 : form.password.length < 10 ? 2 : 3;
  const strengthLabels = ['', 'Weak', 'Good', 'Strong'];
  const strengthColors = ['', 'bg-red-500', 'bg-amber-500', 'bg-green-500'];

  return (
    <div className="min-h-screen flex">
      {/* Left side — Professional dark branding */}
      <div className="hidden lg:flex lg:w-[55%] bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 relative overflow-hidden">
        <div className="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGRlZnM+PHBhdHRlcm4gaWQ9ImciIHdpZHRoPSI2MCIgaGVpZ2h0PSI2MCIgcGF0dGVyblVuaXRzPSJ1c2VyU3BhY2VPblVzZSI+PGNpcmNsZSBjeD0iMzAiIGN5PSIzMCIgcj0iMSIgZmlsbD0icmdiYSgyNTUsMjU1LDI1NSwwLjAzKSIvPjwvcGF0dGVybj48L2RlZnM+PHJlY3QgZmlsbD0idXJsKCNnKSIgd2lkdGg9IjEwMCUiIGhlaWdodD0iMTAwJSIvPjwvc3ZnPg==')] opacity-100" />
        <div className="absolute top-32 right-10 w-72 h-72 bg-accent-500/10 rounded-full blur-[100px]" />
        <div className="absolute bottom-32 left-10 w-96 h-96 bg-primary-500/8 rounded-full blur-[100px]" />

        <div className="relative z-10 flex flex-col justify-between p-12 xl:p-16 w-full">
          <Link to="/" className="flex items-center gap-3 group">
            <div className="w-11 h-11 bg-primary-500 rounded-xl flex items-center justify-center group-hover:scale-105 transition-transform">
              <ShoppingBag className="w-5 h-5 text-white" />
            </div>
            <span className="text-xl font-bold text-white">CloudMart</span>
          </Link>

          <div className="max-w-md">
            <h2 className="text-4xl xl:text-5xl font-bold text-white leading-tight tracking-tight">
              Start your
              <span className="block bg-gradient-to-r from-primary-400 to-accent-400 bg-clip-text text-transparent">
                shopping journey
              </span>
            </h2>
            <p className="mt-5 text-gray-400 text-lg leading-relaxed">
              Create your free account and unlock exclusive benefits, track orders, and enjoy personalized shopping.
            </p>

            <div className="mt-10 grid grid-cols-1 gap-4">
              {[
                { icon: Gift, title: '10% Welcome Discount', desc: 'Automatic discount on your first order' },
                { icon: Heart, title: 'Wishlists & Favorites', desc: 'Save products and get restock alerts' },
                { icon: Star, title: 'Rewards Program', desc: 'Earn points on every purchase' },
              ].map(({ icon: Icon, title, desc }) => (
                <div key={title} className="flex items-start gap-4 p-4 rounded-xl bg-white/5 border border-white/5">
                  <div className="w-10 h-10 rounded-lg bg-primary-500/15 flex items-center justify-center shrink-0">
                    <Icon className="w-5 h-5 text-primary-400" />
                  </div>
                  <div>
                    <p className="text-sm font-semibold text-white">{title}</p>
                    <p className="text-sm text-gray-500 mt-0.5">{desc}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <p className="text-gray-500 text-sm">
            &copy; {new Date().getFullYear()} CloudMart. All rights reserved.
          </p>
        </div>
      </div>

      {/* Right side — Form */}
      <div className="flex-1 flex items-center justify-center px-4 sm:px-8 py-12 bg-gray-50">
        <div className="w-full max-w-[420px] animate-fade-in-up">
          <div className="text-center mb-8 lg:hidden">
            <Link to="/" className="inline-flex items-center gap-2.5">
              <div className="w-11 h-11 bg-gradient-to-br from-primary-500 to-primary-600 rounded-xl flex items-center justify-center shadow-lg shadow-primary-500/20">
                <ShoppingBag className="w-6 h-6 text-white" />
              </div>
              <span className="text-2xl font-bold text-gray-900">CloudMart</span>
            </Link>
          </div>

          <div>
            <div className="mb-8">
              <h2 className="text-2xl sm:text-3xl font-bold text-gray-900 tracking-tight">Create Account</h2>
              <p className="text-gray-500 mt-2">Fill in your details to get started</p>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid grid-cols-2 gap-3">
                <Input
                  label="First Name"
                  type="text"
                  value={form.first_name}
                  onChange={update('first_name')}
                  required
                  placeholder="John"
                  icon={<User className="w-5 h-5" />}
                />
                <Input
                  label="Last Name"
                  type="text"
                  value={form.last_name}
                  onChange={update('last_name')}
                  required
                  placeholder="Doe"
                  icon={<User className="w-5 h-5" />}
                />
              </div>
              <Input
                label="Email address"
                type="email"
                value={form.email}
                onChange={update('email')}
                required
                placeholder="you@example.com"
                icon={<Mail className="w-5 h-5" />}
              />
              <div>
                <Input
                  label="Password"
                  type="password"
                  value={form.password}
                  onChange={update('password')}
                  required
                  placeholder="Minimum 8 characters"
                  icon={<Lock className="w-5 h-5" />}
                />
                {form.password.length > 0 && (
                  <div className="mt-2 flex items-center gap-2">
                    <div className="flex-1 flex gap-1">
                      {[1, 2, 3].map((level) => (
                        <div
                          key={level}
                          className={`h-1.5 flex-1 rounded-full transition-all ${
                            passwordStrength >= level ? strengthColors[passwordStrength] : 'bg-gray-200'
                          }`}
                        />
                      ))}
                    </div>
                    <span className={`text-xs font-medium ${
                      passwordStrength === 1 ? 'text-red-500' : passwordStrength === 2 ? 'text-amber-500' : 'text-green-500'
                    }`}>
                      {strengthLabels[passwordStrength]}
                    </span>
                  </div>
                )}
              </div>

              <div className="pt-2">
                <Button
                  type="submit"
                  loading={loading}
                  className="w-full"
                  size="lg"
                  icon={<ArrowRight className="w-5 h-5" />}
                >
                  Create Account
                </Button>
              </div>
            </form>

            <p className="mt-6 text-center text-xs text-gray-400">
              By creating an account, you agree to our{' '}
              <span className="text-primary-600 font-medium cursor-pointer">Terms of Service</span> and{' '}
              <span className="text-primary-600 font-medium cursor-pointer">Privacy Policy</span>.
            </p>

            <div className="mt-6 text-center">
              <p className="text-sm text-gray-500">
                Already have an account?{' '}
                <Link to="/login" className="text-primary-600 font-semibold hover:text-primary-700 transition-colors">
                  Sign in
                </Link>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
