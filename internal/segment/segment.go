
package segment

import (
    "fmt"
    "path/filepath"
)

func File(
    dir string,
    id int,
) string {

    return filepath.Join(
        dir,
        fmt.Sprintf(
            "%06d.log",
            id,
        ),
    )
}
