package vec3

import "math"

type Vec3 []float32


func (v Vec3) Add(o Vec3) Vec3 {
  return Vec3{v[0] + o[0],
             v[1] + o[1],
             v[2] + o[2]}
}
func (v Vec3) Sub(o Vec3) Vec3 {
  return Vec3{v[0] - o[0],
             v[1] - o[1],
             v[2] - o[2]}
}

func Cross(a Vec3, b Vec3) Vec3 {
  return Vec3{a[1] * b[2] - a[2] * b[1],
             a[2] * b[0] - a[0] * b[2],
             a[0] * b[1] - a[1] * b[0]}
}


func Normalize(a Vec3) Vec3 {
  norm := float32(math.Sqrt(float64(Dot(a, a))))
  return Vec3{a[0]/norm, a[1]/norm, a[2]/norm}
}

func Dot(a Vec3, b Vec3) float32 {
  return a[0] * b[0] + a[1] * b[1] + a[2] * b[2];
}

func Scale(v Vec3, s float32) Vec3 {
  return Vec3{v[0] * s, v[1] * s, v[2] * s}
}

func Add(vecs ...Vec3) Vec3 {
  total := Vec3{0,0,0}
  for  _, v := range vecs {
    total = total.Add(v)
  }
  return total

}


