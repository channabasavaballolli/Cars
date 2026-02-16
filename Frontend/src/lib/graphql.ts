const GRAPHQL_ENDPOINT = "/graphql";

interface GraphQLResponse<T = any> {
  data?: T;
  errors?: { message: string }[];
}

export async function graphqlRequest<T = any>(
  query: string,
  variables?: Record<string, any>,
  token?: string | null
): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(GRAPHQL_ENDPOINT, {
    method: "POST",
    headers,
    body: JSON.stringify({ query, variables }),
  });

  const json: GraphQLResponse<T> = await res.json();

  if (json.errors && json.errors.length > 0) {
    throw new Error(json.errors[0].message);
  }

  return json.data as T;
}

// Queries
export const CARS_QUERY = `
  query {
    cars { id make model year price color mileage }
  }
`;

// Mutations
export const REQUEST_LOGIN = `
  mutation RequestLogin($email: String!) {
    requestLogin(email: $email)
  }
`;

export const VERIFY_LOGIN = `
  mutation VerifyLogin($email: String!, $code: String!) {
    verifyLogin(email: $email, code: $code)
  }
`;

export const CREATE_CAR = `
  mutation CreateCar($make: String!, $model: String!, $year: Int!, $price: Float!, $color: String!, $mileage: Int!) {
    createCar(make: $make, model: $model, year: $year, price: $price, color: $color, mileage: $mileage) { id }
  }
`;

export const UPDATE_CAR = `
  mutation UpdateCar($id: Int!, $make: String, $model: String, $year: Int, $price: Float, $color: String, $mileage: Int) {
    updateCar(id: $id, make: $make, model: $model, year: $year, price: $price, color: $color, mileage: $mileage) { id }
  }
`;

export const DELETE_CAR = `
  mutation DeleteCar($id: Int!) {
    deleteCar(id: $id)
  }
`;
