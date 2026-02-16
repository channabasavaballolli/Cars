import { useEffect, useState, useCallback } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { graphqlRequest, CARS_QUERY, CREATE_CAR, UPDATE_CAR, DELETE_CAR } from "@/lib/graphql";
import Navbar from "@/components/Navbar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/components/ui/table";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { AlertDialog, AlertDialogContent, AlertDialogHeader, AlertDialogTitle, AlertDialogDescription, AlertDialogFooter, AlertDialogAction, AlertDialogCancel } from "@/components/ui/alert-dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";
import { Plus, Pencil, Trash2, Loader2, Car } from "lucide-react";
import { motion } from "framer-motion";

interface CarData {
  id: number;
  make: string;
  model: string;
  year: number;
  price: number;
  color: string;
  mileage: number;
}

const emptyForm = { make: "", model: "", year: "", price: "", color: "", mileage: "" };

function formatPrice(price: number) {
  return new Intl.NumberFormat("en-US", { style: "currency", currency: "USD", maximumFractionDigits: 0 }).format(price);
}

export default function Dashboard() {
  const { token } = useAuth();
  const [cars, setCars] = useState<CarData[]>([]);
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [editingCar, setEditingCar] = useState<CarData | null>(null);
  const [form, setForm] = useState(emptyForm);
  const [saving, setSaving] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState<CarData | null>(null);
  const [deleting, setDeleting] = useState(false);

  const fetchCars = useCallback(() => {
    setLoading(true);
    graphqlRequest<{ cars: CarData[] }>(CARS_QUERY)
      .then((d) => setCars(d.cars))
      .catch((e) => toast.error(e.message))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => { fetchCars(); }, [fetchCars]);

  const openAdd = () => {
    setEditingCar(null);
    setForm(emptyForm);
    setModalOpen(true);
  };

  const openEdit = (car: CarData) => {
    setEditingCar(car);
    setForm({
      make: car.make,
      model: car.model,
      year: String(car.year),
      price: String(car.price),
      color: car.color,
      mileage: String(car.mileage),
    });
    setModalOpen(true);
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      const vars = {
        make: form.make,
        model: form.model,
        year: parseInt(form.year),
        price: parseFloat(form.price),
        color: form.color,
        mileage: parseInt(form.mileage),
      };
      if (editingCar) {
        await graphqlRequest(UPDATE_CAR, { id: editingCar.id, ...vars }, token);
        toast.success("Vehicle updated");
      } else {
        await graphqlRequest(CREATE_CAR, vars, token);
        toast.success("Vehicle added");
      }
      setModalOpen(false);
      fetchCars();
    } catch (err: any) {
      toast.error(err.message);
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    if (!deleteTarget) return;
    setDeleting(true);
    try {
      await graphqlRequest(DELETE_CAR, { id: deleteTarget.id }, token);
      toast.success("Vehicle deleted");
      setDeleteTarget(null);
      fetchCars();
    } catch (err: any) {
      toast.error(err.message);
    } finally {
      setDeleting(false);
    }
  };

  const updateField = (field: string, value: string) => setForm((p) => ({ ...p, [field]: value }));

  return (
    <div className="min-h-screen gradient-hero">
      <Navbar />

      <div className="container mx-auto pt-24 pb-12 px-4">
        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Dashboard</h1>
            <p className="text-sm text-muted-foreground mt-1">{cars.length} vehicles in inventory</p>
          </div>
          <Button onClick={openAdd} className="gap-2 font-semibold">
            <Plus className="h-4 w-4" /> Add Vehicle
          </Button>
        </motion.div>

        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} transition={{ delay: 0.1 }} className="glass rounded-xl overflow-hidden">
          {loading ? (
            <div className="p-6 space-y-3">
              {Array.from({ length: 5 }).map((_, i) => <Skeleton key={i} className="h-12 w-full rounded-lg" />)}
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow className="border-border/50 hover:bg-transparent">
                  <TableHead className="text-muted-foreground">Make</TableHead>
                  <TableHead className="text-muted-foreground">Model</TableHead>
                  <TableHead className="text-muted-foreground">Year</TableHead>
                  <TableHead className="text-muted-foreground">Price</TableHead>
                  <TableHead className="text-muted-foreground">Mileage</TableHead>
                  <TableHead className="text-muted-foreground">Color</TableHead>
                  <TableHead className="text-muted-foreground text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {cars.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center py-12 text-muted-foreground">
                      <Car className="h-8 w-8 mx-auto mb-2 opacity-40" />
                      No vehicles yet. Add one to get started.
                    </TableCell>
                  </TableRow>
                ) : (
                  cars.map((car) => (
                    <TableRow key={car.id} className="border-border/30 hover:bg-secondary/30">
                      <TableCell className="font-medium">{car.make}</TableCell>
                      <TableCell>{car.model}</TableCell>
                      <TableCell>{car.year}</TableCell>
                      <TableCell className="font-semibold text-primary">{formatPrice(car.price)}</TableCell>
                      <TableCell>{car.mileage?.toLocaleString()} mi</TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <div className="w-3 h-3 rounded-full border border-border" style={{ backgroundColor: car.color?.toLowerCase() }} />
                          <span className="capitalize">{car.color}</span>
                        </div>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex justify-end gap-1">
                          <Button variant="ghost" size="icon" onClick={() => openEdit(car)}>
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button variant="ghost" size="icon" className="text-destructive hover:text-destructive" onClick={() => setDeleteTarget(car)}>
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          )}
        </motion.div>
      </div>

      {/* Add/Edit Modal */}
      <Dialog open={modalOpen} onOpenChange={setModalOpen}>
        <DialogContent className="glass-strong sm:max-w-md">
          <DialogHeader>
            <DialogTitle>{editingCar ? "Edit Vehicle" : "Add Vehicle"}</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSave} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              {[
                { label: "Make", field: "make", type: "text", placeholder: "Toyota" },
                { label: "Model", field: "model", type: "text", placeholder: "Camry" },
                { label: "Year", field: "year", type: "number", placeholder: "2024" },
                { label: "Price", field: "price", type: "number", placeholder: "25000" },
                { label: "Color", field: "color", type: "text", placeholder: "White" },
                { label: "Mileage", field: "mileage", type: "number", placeholder: "0" },
              ].map(({ label, field, type, placeholder }) => (
                <div key={field} className="space-y-1.5">
                  <Label className="text-xs text-muted-foreground">{label}</Label>
                  <Input
                    type={type}
                    placeholder={placeholder}
                    value={(form as any)[field]}
                    onChange={(e) => updateField(field, e.target.value)}
                    className="bg-secondary/50 border-border/50"
                    required
                  />
                </div>
              ))}
            </div>
            <DialogFooter>
              <Button type="submit" disabled={saving} className="w-full font-semibold">
                {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : editingCar ? "Update" : "Add Vehicle"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation */}
      <AlertDialog open={!!deleteTarget} onOpenChange={(open) => !open && setDeleteTarget(null)}>
        <AlertDialogContent className="glass-strong">
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Vehicle?</AlertDialogTitle>
            <AlertDialogDescription>
              This will permanently remove the {deleteTarget?.year} {deleteTarget?.make} {deleteTarget?.model} from inventory.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={handleDelete} disabled={deleting} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
              {deleting ? <Loader2 className="h-4 w-4 animate-spin" /> : "Delete"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
