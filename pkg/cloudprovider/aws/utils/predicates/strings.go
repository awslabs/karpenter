/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package predicates

// WithinStrings returns a func that returns true if string is within strings.
func WithinStrings(allowed []string) func(string) bool {
	return func(actual string) bool {
		for _, expected := range allowed {
			if expected == actual {
				return true
			}
		}
		return false
	}
}
