package fortune

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
)

type FortuneTeller struct {
}
type SuggestionRequest struct {
}
type SuggestionResponse struct {
}
type Server struct {
}

// Adding information to errors
func openFile() error {
	// 보기 좋게 지역변수로 선언한 것뿐
	var _ *os.File
	var err error

	// Good:
	// When adding information to errors, avoid redundant information that the underlying error already provides.
	// The os package, for instance, already includes path information in its errors.
	if _, err = os.Open("settings.txt"); err != nil {
		return fmt.Errorf("launch codes unavailable: %v", err)
	}
	// Output:
	//
	// launch codes unavailable: open settings.txt: no such file or directory

	// Bad:
	if _, err = os.Open("settings.txt"); err != nil {
		return fmt.Errorf("could not open settings.txt: %v", err)
	}
	// Output:
	//
	// could not open settings.txt: open settings.txt: no such file or directory

	// 당연히 아래처럼 failed: 같은 것만 덧붙이는 것은 노쓸모
	// Bad:
	return fmt.Errorf("failed: %v", err) // just return err instead
}

// when wrapping errors %v vs. %w

// 1. %v for simple annotation or new error
// 에러를 문자열로만 매핑
// 시스템 밖으로 나갈 때는 내부 에러를 그대로 내보내지 말고, 외부에서 이해 가능한 표준 에러로 “번역”하라.
// Good:
func (*FortuneTeller) SuggestFortune(ctx context.Context, _ SuggestionRequest) (SuggestionResponse, error) {
	var err error
	var response SuggestionResponse
	if err != nil {
		return response, fmt.Errorf("couldn't find fortune database: %v", err)
	}

	if err != nil {
		// Or use fmt.Errorf with the %w verb if deliberately wrapping an
		// error which the caller is meant to unwrap.
		return response, status.Errorf(codes.Internal, "couldn't find fortune database")
	}
	return response, nil
}

// 2. %w (wrap) for programmatic inspection and error chaining
// 에러 랩핑에 사용
// 추가적인 컨텍스트 (실패했을 때 어떤 동작이 수행됐는지에 대한 정보 등)와 함께 에러를 풍부하게 하고 싶은 경우 사용
// errors.Is / errors.As 가능
// Good:
func (s Server) internalFunction(ctx context.Context) error {
	var err error
	if err != nil {
		return fmt.Errorf("couldn't find remote file: %w", err)
	}
	return nil
}
