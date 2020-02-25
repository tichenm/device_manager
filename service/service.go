package service

//
import(
	//"zhiyuan/ai_dormitory_apis/school_affairs/conf"
	"zhiyuan/scaffold/internal/dao"
	"context"
	"zhiyuan/scaffold/configs"
)
//// Service service.
type Service struct {
	//	ac  *paladin.Map
	//	dao *dao.Dao
	c           *configs.Config
	dao         *dao.Dao
}
//
// New new a service and return.
func New(c *configs.Config) (s *Service) {
	//var ac = new(paladin.TOML)
	//if err := paladin.Watch("application.toml", ac); err != nil {
	//	panic(err)
	//}
	s = &Service{
		c:  c,
		dao: dao.New(c),
	}
	return s
}
//
// Ping ping the resource.
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}
//
// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}

