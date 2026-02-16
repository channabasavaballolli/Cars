import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { graphqlRequest, CARS_QUERY } from "@/lib/graphql";
import { Skeleton } from "@/components/ui/skeleton";
import { Car, Gauge, Calendar, DollarSign } from "lucide-react";
import Navbar from "@/components/Navbar";

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

function CarCard({ car, index }: { car: CarData; index: number }) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 30 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4, delay: index * 0.06 }}
      className="group glass rounded-xl overflow-hidden hover:border-primary/40 transition-all duration-300"
    >
      {/* Color accent bar */}
      <div className="h-1.5 w-full" style={{ backgroundColor: car.color?.toLowerCase() || "hsl(210,100%,56%)" }} />

      <div className="p-5 space-y-4">
        <div className="flex items-start justify-between">
          <div>
            <h3 className="text-lg font-bold text-foreground">{car.make}</h3>
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
        <div className="opacity-0 group-hover:opacity-100 transition-opacity duration-300 pt-2 border-t border-border/50">
          <p className="text-xs text-muted-foreground">
            Color: <span className="text-foreground capitalize">{car.color || "N/A"}</span>
          </p>
        </div>
      </div>
    </motion.div>
  );
}

export default function Index() {
  const [cars, setCars] = useState<CarData[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    graphqlRequest<{ cars: CarData[] }>(CARS_QUERY)
      .then((data) => setCars(data.cars))
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

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
          {loading ? (
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
                <CarCard key={car.id} car={car} index={i} />
              ))}
            </div>
          )}
        </div>
      </section>
    </div>
  );
}
