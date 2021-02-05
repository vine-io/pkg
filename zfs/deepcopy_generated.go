// +build !ignore_autogenerated

// Copyright 2021 lack
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Code generated by deepcopy-gen. Do NOT EDIT.

package zfs

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *Pool) DeepCopyInto(out *Pool) {
	*out = *in
	out.Types = in.Types
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make(map[string]*Volume, len(*in))
		for key, val := range *in {
			var outVal *Volume
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(Volume)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
	if in.Snapshots != nil {
		in, out := &in.Snapshots, &out.Snapshots
		*out = make(map[string]*Snapshot, len(*in))
		for key, val := range *in {
			var outVal *Snapshot
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(Snapshot)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
	if in.Clones != nil {
		in, out := &in.Clones, &out.Clones
		*out = make(map[string]*Volume, len(*in))
		for key, val := range *in {
			var outVal *Volume
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(Volume)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
	return
}

// DeepCopy is an auto-generated deepcopy function, copying the receiver, creating a new Pool.
func (in *Pool) DeepCopy() *Pool {
	if in == nil {
		return nil
	}
	out := new(Pool)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *Snapshot) DeepCopyInto(out *Snapshot) {
	*out = *in
	out.Types = in.Types
	if in.Clones != nil {
		in, out := &in.Clones, &out.Clones
		*out = make(map[string]*Volume, len(*in))
		for key, val := range *in {
			var outVal *Volume
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(Volume)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
	return
}

// DeepCopy is an auto-generated deepcopy function, copying the receiver, creating a new Snapshot.
func (in *Snapshot) DeepCopy() *Snapshot {
	if in == nil {
		return nil
	}
	out := new(Snapshot)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *Volume) DeepCopyInto(out *Volume) {
	*out = *in
	out.Types = in.Types
	if in.Snapshots != nil {
		in, out := &in.Snapshots, &out.Snapshots
		*out = make(map[string]*Snapshot, len(*in))
		for key, val := range *in {
			var outVal *Snapshot
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(Snapshot)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
	return
}

// DeepCopy is an auto-generated deepcopy function, copying the receiver, creating a new Volume.
func (in *Volume) DeepCopy() *Volume {
	if in == nil {
		return nil
	}
	out := new(Volume)
	in.DeepCopyInto(out)
	return out
}
