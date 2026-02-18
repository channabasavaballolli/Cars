import { useEffect, useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { motion, AnimatePresence } from "framer-motion";
import { graphqlRequest, CARS_QUERY } from "@/lib/graphql";
import { Skeleton } from "@/components/ui/skeleton";
import { Car, Gauge, Calendar, DollarSign, X } from "lucide-react";
import Navbar from "@/components/Navbar";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";

interface CarData {
  id: number;
  make: string;
  model: string;
  year: number;
  price: number;
  color: string;
  mileage: number;
}

function formatPrice(price: number) {
  return new Intl.NumberFormat("en-US", { style: "currency", currency: "USD", maximumFractionDigits: 0 }).format(price);
}

function CarCard({ car, index, onClick }: { car: CarData; index: number; onClick: () => void }) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 30 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4, delay: index * 0.06 }}
      whileHover={{ y: -5, transition: { duration: 0.2 } }}
      onClick={onClick}
      className="group glass rounded-xl overflow-hidden hover:border-primary/40 transition-all duration-300 cursor-pointer"
    >
      {/* Color accent bar */}
      <div className="h-1.5 w-full" style={{ backgroundColor: car.color?.toLowerCase() || "hsl(210,100%,56%)" }} />

      <div className="p-5 space-y-4">
        <div className="flex items-start justify-between">
          <div>
            <h3 className="text-lg font-bold text-foreground group-hover:text-primary transition-colors">{car.make}</h3>
            <p className="text-sm text-muted-foreground">{car.model}</p>
          </div>
          <span className="text-xs font-medium px-2.5 py-1 rounded-full bg-primary/10 text-primary border border-primary/20">
            {car.year}
          </span>
        </div>

        <div className="grid grid-cols-2 gap-3 text-sm">
          <div className="flex items-center gap-2 text-muted-foreground">
            <DollarSign className="h-3.5 w-3.5 text-accent" />
            <span className="text-foreground font-semibold">{formatPrice(car.price)}</span>
          </div>
          <div className="flex items-center gap-2 text-muted-foreground">
            <Gauge className="h-3.5 w-3.5 text-accent" />
            <span>{car.mileage?.toLocaleString()} mi</span>
          </div>
        </div>

        {/* Hover overlay */}
        <div className="opacity-0 group-hover:opacity-100 transition-opacity duration-300 pt-2 border-t border-border/50 flex justify-between items-center">
          <p className="text-xs text-muted-foreground">
            Color: <span className="text-foreground capitalize">{car.color || "N/A"}</span>
          </p>
          <span className="text-xs text-primary font-medium">View Details →</span>
        </div>
      </div>
    </motion.div>
  );
}

export default function Index() {
  const { isAuthenticated } = useAuth();
  const [cars, setCars] = useState<CarData[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedCar, setSelectedCar] = useState<CarData | null>(null);

  useEffect(() => {
    if (isAuthenticated) {
      setLoading(true);
      graphqlRequest<{ cars: CarData[] }>(CARS_QUERY)
        .then((data) => setCars(data.cars))
        .catch((e) => setError(e.message))
        .finally(() => setLoading(false));
    }
  }, [isAuthenticated]);

  return (
    <div className="min-h-screen gradient-hero">
      <Navbar />

      {/* Hero */}
      <section className="pt-32 pb-16 px-4">
        <div className="container mx-auto text-center max-w-3xl">
          <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.6 }}>
            <span className="inline-block text-xs font-semibold tracking-widest uppercase text-accent mb-4">
              Curated Collection
            </span>
            <h1 className="text-5xl md:text-6xl font-black tracking-tight text-foreground mb-4">
              Premium <span className="text-primary">Inventory</span>
            </h1>
            <p className="text-lg text-muted-foreground max-w-xl mx-auto">
              Explore our handpicked selection of premium vehicles. Quality meets elegance.
            </p>
          </motion.div>
        </div>
      </section>

      {/* Grid */}
      <section className="pb-24 px-4">
        <div className="container mx-auto">
          {!isAuthenticated ? (
            <div className="text-center py-20 bg-secondary/20 rounded-2xl border border-white/10 backdrop-blur-sm">
              <Car className="h-16 w-16 mx-auto text-primary mb-6 opacity-80" />
              <h2 className="text-2xl font-bold text-white mb-4">Exclusive Inventory</h2>
              <p className="text-muted-foreground max-w-md mx-auto mb-8">
                Our premium collection is available exclusively to registered members.
                Please login to view our available vehicles.
              </p>
              <div className="flex justify-center gap-4">
                <a href="/login" className="px-6 py-2 bg-primary text-primary-foreground font-semibold rounded-lg hover:bg-primary/90 transition-colors">
                  Member Login
                </a>
                <a href="/admin-login" className="px-6 py-2 bg-secondary text-secondary-foreground font-semibold rounded-lg hover:bg-secondary/80 transition-colors">
                  Admin Portal
                </a>
              </div>
            </div>
          ) : loading ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
              {Array.from({ length: 8 }).map((_, i) => (
                <Skeleton key={i} className="h-48 rounded-xl" />
              ))}
            </div>
          ) : error ? (
            <div className="text-center py-20">
              <p className="text-destructive mb-2">Failed to load inventory</p>
              <p className="text-sm text-muted-foreground">{error}</p>
            </div>
          ) : cars.length === 0 ? (
            <div className="text-center py-20">
              <Car className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
              <p className="text-muted-foreground">No vehicles in inventory yet.</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
              {cars.map((car, i) => (
                <CarCard
                  key={car.id}
                  car={car}
                  index={i}
                  onClick={() => setSelectedCar(car)}
                />
              ))}
            </div>
          )}
        </div>
      </section>

      {/* Car Details Modal */}
      <Dialog open={!!selectedCar} onOpenChange={(open) => !open && setSelectedCar(null)}>
        <DialogContent className="glass-strong border-primary/20 max-w-md sm:max-w-lg p-0 overflow-hidden gap-0">
          {selectedCar && (
            <>
              <div className="absolute top-0 left-0 w-full h-1.5" style={{ backgroundColor: selectedCar.color?.toLowerCase() || "hsl(210,100%,56%)" }} />

              <DialogHeader className="p-6 pb-2 text-left">
                <div className="flex justify-between items-start">
                  <div>
                    <span className="text-xs font-medium text-primary bg-primary/10 px-2 py-0.5 rounded border border-primary/20 mb-2 inline-block">
                      {selectedCar.year} Model
                    </span>
                    <DialogTitle className="text-3xl font-bold tracking-tight mt-1">{selectedCar.make}</DialogTitle>
                    <DialogDescription className="text-lg font-medium text-foreground/80">{selectedCar.model}</DialogDescription>
                  </div>
                </div>
              </DialogHeader>

              <div className="p-6 pt-2 space-y-6">
                <div className="grid grid-cols-2 gap-4">
                  <div className="bg-secondary/30 p-4 rounded-xl border border-border/50">
                    <div className="flex items-center gap-2 text-muted-foreground mb-1">
                      <DollarSign className="h-4 w-4 text-accent" />
                      <span className="text-xs uppercase tracking-wider font-semibold">Price</span>
                    </div>
                    <p className="text-2xl font-bold text-foreground">{formatPrice(selectedCar.price)}</p>
                  </div>
                  <div className="bg-secondary/30 p-4 rounded-xl border border-border/50">
                    <div className="flex items-center gap-2 text-muted-foreground mb-1">
                      <Gauge className="h-4 w-4 text-accent" />
                      <span className="text-xs uppercase tracking-wider font-semibold">Mileage</span>
                    </div>
                    <p className="text-2xl font-bold text-foreground">{selectedCar.mileage.toLocaleString()} mi</p>
                  </div>
                </div>

                <div className="space-y-3">
                  <h4 className="text-sm font-semibold text-muted-foreground uppercase tracking-widest border-b border-border/50 pb-2">
                    Vehicle Specifications
                  </h4>
                  <div className="grid grid-cols-2 gap-y-3 text-sm">
                    <div className="flex justify-between pr-4">
                      <span className="text-muted-foreground">Make</span>
                      <span className="font-medium">{selectedCar.make}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Year</span>
                      <span className="font-medium">{selectedCar.year}</span>
                    </div>
                    <div className="flex justify-between pr-4">
                      <span className="text-muted-foreground">Model</span>
                      <span className="font-medium">{selectedCar.model}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Color</span>
                      <div className="flex items-center gap-2">
                        <div className="w-3 h-3 rounded-full border border-border" style={{ backgroundColor: selectedCar.color?.toLowerCase() }} />
                        <span className="font-medium capitalize">{selectedCar.color}</span>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="pt-2">
                  <button
                    onClick={() => setSelectedCar(null)}
                    className="w-full py-3 rounded-xl bg-primary text-primary-foreground font-semibold shadow-lg shadow-primary/20 hover:bg-primary/90 transition-all hover:scale-[1.02] active:scale-[0.98]"
                  >
                    Close Details
                  </button>
                </div>
              </div>
            </>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
