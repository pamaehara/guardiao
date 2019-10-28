package api

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	YYYY_MM_DD_T_HH_MM_SS = "2006-01-02T15:04:05"
	SHORT_DATE            = "01-02"
	SHORT_TIME            = "15:04"
	path                  = "/gwaereo/v0/flights"
	key                   = "Y5a7r8e1q2w49EH0"
	LOG_LEVEL_SPACE       = "     "
)

type IataValue struct {
	Iata string `json:"iata"`
	Name string `json:"name"`
}

type Header struct {
	Name  string
	Value string
}

type RequestApi struct {
	UrlGwapi        string `json:"url_gwapi"`
	UrlPrecificador string `json:"url_precificador"`
	Url             string `json:"url"`
	UrlNotify       string `json:"urlNotify"`
	Credential      `json:"credential"`
}

type Credential struct {
	User      string `json:"user"`
	Passwd    string `json:"passwd"`
	BranchID  string `json:"branchId"`
	AgencyID  string `json:"agencyId"`
	GroupID   string `json:"groupId"`
	AgentSign string `json:"agentSign"`
}

func (r RequestApi) GetHeaders() []Header {
	return []Header{
		Header{
			Name:  "Gtw-Username",
			Value: r.Credential.User,
		},
		Header{
			Name:  "Gtw-Password",
			Value: r.Credential.Passwd,
		},
		Header{
			Name:  "Gtw-Branch-Id",
			Value: r.Credential.BranchID,
		},
		Header{
			Name:  "Gtw-Agent-Sign",
			Value: r.Credential.AgentSign,
		},
		Header{
			Name:  "Gtw-Agency-Id",
			Value: r.Credential.AgencyID,
		},
		Header{
			Name:  "Gtw-Group-Id",
			Value: r.Credential.GroupID,
		},
		Header{
			Name:  "Accept",
			Value: "application/json; charset=utf-8",
		},
		Header{
			Name:  "Content-Type",
			Value: "application/json; charset=utf-8",
		},
		Header{
			Name:  "Accept-Charset",
			Value: "utf-8",
		},
	}
}

func (r RequestApi) GetSoapHeaders() []Header {
	return []Header{
		Header{
			Name:  "Content-Type",
			Value: "Content-Type: text/xml;charset=UTF-8",
		},
		Header{
			Name:  "SOAPAction",
			Value: "",
		},
	}
}

func (a RequestApi) GetURL() string {
	return a.Url + path
}

func (a RequestApi) GetURLV1() string {
	return strings.Replace(a.Url+path, "v0", "v1", -1)
}

type RateToken struct {
	Token   string
	decoded string
}

func GetMarkupDescription(markupID string, request RequestApi) MarkupDescription {
	client := http.DefaultClient

	jwtPrec := GetJWTPrec(request.UrlPrecificador, request.Credential.User, request.Credential.Passwd, request.Credential.BranchID, request.Credential.AgencyID)

	req, err := http.NewRequest("GET", request.UrlPrecificador+"/rs/v1/configuracoesPrecificacao/"+markupID, nil)
	req.Header.Add("Authorization", "Bearer "+jwtPrec)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		log.Info(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Info(err)
	}

	response := MarkupDescription{}
	json.Unmarshal(body, &response)

	return response
}

func (r *RateToken) getDecoded() string {
	if r.decoded == "" {
		token := r.Token
		if decoded, err := decryptECB(token); err == nil {
			r.decoded = decoded
			return r.decoded
		}
		b, _ := b64.URLEncoding.DecodeString(token)
		r.decoded = string(b)
	}
	return r.decoded
}

func decipherECB(input []byte) (string, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	ecb := newECBDecrypter(c)

	dst := make([]byte, len(input))

	err = ecb.CryptBlocks(dst, input)
	if err != nil {
		return "", err
	}

	decoded := string(dst)
	return decoded, nil
}

func decryptECB(token string) (string, error) {
	b, _ := b64.URLEncoding.DecodeString(token)

	decoded, err := decipherECB(b)

	if err != nil {
		return "", err
	}

	b642bytes, _ := b64.URLEncoding.DecodeString(decoded)

	gz, err := gzip.NewReader(bytes.NewReader(b642bytes))
	if err != nil {
		//	log.Printf("erro no gzip %v\n", err)
		return "", err
	}
	defer gz.Close()

	s, err := ioutil.ReadAll(gz)
	if err != nil {
		//	log.Printf("erro no readall %v\n", err)
		return "", err
	}
	s, err = b64.URLEncoding.DecodeString(string(s))
	if err != nil {
		//	log.Printf("erro no terceiro b64 %v\n", err)
		return "", err
	}

	return string(s), nil
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbDecrypter ecb

func newECBDecrypter(b cipher.Block) *ecbDecrypter {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) error {
	if len(src)%x.blockSize != 0 {
		return fmt.Errorf("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		return fmt.Errorf("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
	return nil
}

// <rateToken bri="1000" cmi="AZUL" cid="162" cur="BRL" ect="BR" dtf="2018-09-09T13:45:00.000Z" est="1076" ezi="1076" mkp="1.0000" ofd="" pxs="30" pkg="VHI" ppy="AD" plc="BRL" ftp="M" pot="1005.80" pwt="1063.55" prd="AIR" sdt="2018-07-15T12:18:25.312Z" sgk="" sot="1005.80" swt="1063.55" sct="BR" dti="2018-09-09T09:20:00.000Z" sst="SP" szi="8983" mki="5976" prf="0.00" mec="CIA" cpp="0" rtk="F+"/>
//<rateToken bri="1000" cmi="AMD" cid="33" cur="BRL" ect="US" dtf="2018-09-09T13:14:00.000Z" est="13528" ezi="13528" mkp="0.7500" ofd="{g:'AMD_r:0',a:'PR',b:'QSEB22106',c:'QSEB22106',d:'QSEB22106',e:[5821],f:[{a:'ADT',b:1,c:1388.84,d:366.07,e:0.00,h:0.00,f:[&quot;PFA&quot;],g:[30]}],i:'bebecbac7eb741be'}" pxs="30" pkg="STANDALONE" ppy="CM" plc="BRL" ftp="C" pot="1388.84" pwt="1754.91" prm="379987" prd="AIR" rtc="379987" sdt="2018-07-17T07:08:35.395Z" sgk="{c:[{a:'CM702-0909180300',b:'GRUPTY0909180801OM',d:'738',e:'OAAAQZ2P'},{a:'CM334-0909180908',b:'PTYMIA0909181314OM',d:'738',e:'OAAAQZ2P'}]}" sot="1851.79" swt="2254.90" sct="BR" dti="2018-09-09T03:00:00.000Z" sst="SP" szi="9141" mki="5821" prf="0.00" mec="OWN" cpp="0"/>
//

type MarkupDescription struct {
	ID                  int         `json:"id"`
	Nome                string      `json:"nome"`
	Prioridade          int         `json:"prioridade"`
	ComissaoFornecedor  float64     `json:"comissaoFornecedor"`
	IncentivoFornecedor float64     `json:"incentivoFornecedor"`
	Tourcode            string      `json:"tourcode"`
	Endosso             string      `json:"endosso"`
	ForcarEmissao       bool        `json:"forcarEmissao"`
	ZerarRepasse        bool        `json:"zerarRepasse"`
	MensagemTriagem     interface{} `json:"mensagemTriagem"`
	Versao              int         `json:"versao"`
	IDOrigem            int         `json:"idOrigem"`
	IDConfigEditavel    int         `json:"idConfigEditavel"`
	Editavel            bool        `json:"editavel"`
	Ativo               bool        `json:"ativo"`
	CiasAereas          []struct {
		Codigo      string `json:"codigo"`
		ID          int    `json:"id"`
		Nome        string `json:"nome"`
		Ativo       bool   `json:"ativo"`
		RegraTaxaDu struct {
			InformacaoNacional struct {
				ValorMinimo   int    `json:"valorMinimo"`
				Porcentagem   int    `json:"porcentagem"`
				Moeda         string `json:"moeda"`
				CobraTaxaBebe bool   `json:"cobraTaxaBebe"`
			} `json:"informacaoNacional"`
			InformacaoInternacional struct {
				ValorMinimo   int    `json:"valorMinimo"`
				Porcentagem   int    `json:"porcentagem"`
				Moeda         string `json:"moeda"`
				CobraTaxaBebe bool   `json:"cobraTaxaBebe"`
			} `json:"informacaoInternacional"`
			Ativo bool `json:"ativo"`
		} `json:"regraTaxaDu"`
		NomeTratado string `json:"nomeTratado"`
	} `json:"ciasAereas"`
	StatusPublicacao       string        `json:"statusPublicacao"`
	TipoConfigPrecificacao string        `json:"tipoConfigPrecificacao"`
	Excecoes               []interface{} `json:"excecoes"`
	Markups                []struct {
		ID      int `json:"id"`
		Produto struct {
			ID        int    `json:"id"`
			Nome      string `json:"nome"`
			Descricao string `json:"descricao"`
			Ativo     bool   `json:"ativo"`
		} `json:"produto"`
		ConfiguracaoPrecificacaoID int     `json:"configuracaoPrecificacaoId"`
		MarkupBusca                float64 `json:"markupBusca"`
		Fee                        float64 `json:"fee"`
		ComissaoCliente            float64 `json:"comissaoCliente"`
		IncentivoCliente           float64 `json:"incentivoCliente"`
		RepasseCliente             float64 `json:"repasseCliente"`
		ZeraDu                     bool    `json:"zeraDu"`
		TaxaBolsa                  float64 `json:"taxaBolsa"`
		FatoresMap                 struct {
			MARKUP           float64 `json:"MARKUP"`
			FEE              float64 `json:"FEE"`
			COMISSAOCLIENTE  float64 `json:"COMISSAO_CLIENTE"`
			REPASSECLIENTE   float64 `json:"REPASSE_CLIENTE"`
			INCENTIVOCLIENTE float64 `json:"INCENTIVO_CLIENTE"`
		} `json:"fatoresMap"`
		MarkupEmFatorDivisao float64 `json:"markupEmFatorDivisao"`
	} `json:"markups"`
	TipoMarkup      string      `json:"tipoMarkup"`
	AcordoComercial interface{} `json:"acordoComercial"`
	NacInt          string      `json:"nacInt"`
	Empresas        []struct {
		ID         int    `json:"id"`
		Referencia string `json:"referencia"`
		Nome       string `json:"nome"`
		Ativo      bool   `json:"ativo"`
		Perfil     string `json:"perfil"`
		Refencia   string `json:"refencia"`
	} `json:"empresas"`
	Restricoes []struct {
		ID                int         `json:"id"`
		RestricaoMaeID    interface{} `json:"restricaoMaeId"`
		Nome              string      `json:"nome"`
		Valor             string      `json:"valor"`
		TipoOperador      string      `json:"tipoOperador"`
		TipoRestricao     string      `json:"tipoRestricao"`
		TipoAgrupamento   string      `json:"tipoAgrupamento"`
		TipoDadoRestricao struct {
			Tipo          string `json:"tipo"`
			Campo         string `json:"campo"`
			Periodo       bool   `json:"periodo"`
			TipoRestricao string `json:"tipoRestricao"`
		} `json:"tipoDadoRestricao"`
		Restricoes []interface{} `json:"restricoes"`
		Valores    []struct {
			Nome          string `json:"nome"`
			Valor         string `json:"valor"`
			TipoOperador  string `json:"tipoOperador"`
			ID            int    `json:"id"`
			TipoRestricao string `json:"tipoRestricao"`
		} `json:"valores"`
	} `json:"restricoes"`
	Historico []interface{} `json:"historico"`
}

func (description MarkupDescription) Print() string {
	var message []string
	message = append(message, fmt.Sprintln("Markup Description: "))
	message = append(message, fmt.Sprintln("\t Id: ", description.ID))
	message = append(message, fmt.Sprintln("\t Nome: ", description.Nome))
	message = append(message, fmt.Sprintln("\t Tipo: ", description.TipoMarkup))
	message = append(message, fmt.Sprintln(fmt.Sprintf("\t Comissão | Incentivo: %.2f | %.2f \n", description.ComissaoFornecedor, description.IncentivoFornecedor)))
	message = append(message, fmt.Sprintln(fmt.Sprintf("\t Endosso | Tourcode: %s | %s \n", description.Endosso, description.Tourcode)))
	return strings.Join(message, "")
}

func GetCredential(user string) Credential {
	switch {
	case user == "cvc" || user == "WEB":
		return Credential{
			User:      "scvc",
			Passwd:    "scvc",
			BranchID:  "1000",
			AgentSign: "WEB",
		}
	case user == "lvl" || user == "LIV":
		return Credential{
			User:      "livelo",
			Passwd:    "8u4abux6_@UgAstA",
			BranchID:  "1003",
			AgentSign: "LVL",
		}
	case user == "spp":
		return Credential{
			User:      "app_sub",
			Passwd:    "F@987654",
			BranchID:  "1020",
			AgentSign: "APP",
		}
	case user == "sv" || user == "SUB":
		return Credential{
			User:      "sub_contr",
			Passwd:    "consrah2015",
			BranchID:  "1020",
			AgentSign: "SUB",
		}
	case user == "lojas" || user == "LOJ":
		return Credential{
			User:      "lojascvc",
			Passwd:    "lojascvc",
			BranchID:  "9000",
			AgentSign: "LOJ",
		}
	case user == "vm":
		return Credential{
			User:      "Lojinhavm",
			Passwd:    "lojinhasVM",
			BranchID:  "9000",
			AgentSign: "LOJ",
		}
	case user == "ra" || user == "RA":
		return Credential{
			User:      "integracaora",
			Passwd:    "integracaora",
			BranchID:  "215",
			AgencyID:  "188",
			AgentSign: "RA",
		}
	case user == "mm":
		return Credential{
			User:      "maxmilhas",
			Passwd:    "CH4bDuPP3m",
			BranchID:  "1000",
			AgencyID:  "249046",
			AgentSign: "MM",
		}

	case user == "esf":
		return Credential{
			User:      "Esfera Tur",
			Passwd:    "1234",
			BranchID:  "5984",
			AgentSign: "ESF",
		}
		//89247 -> grupo no RF
	case user == "bct":
		return Credential{
			User:      "brasilctgw",
			Passwd:    "br1912ct",
			BranchID:  "215",
			AgencyID:  "3861",
			AgentSign: "BCT",
		}
	case user == "vs":
		return Credential{
			User:      "VISUAL",
			Passwd:    "Visual@@",
			BranchID:  "1052",
			AgentSign: "VIS",
		}
	case user == "tr":
		return Credential{
			User:      "TREND",
			Passwd:    "123456",
			BranchID:  "1040",
			AgentSign: "TRE",
		}
	case user == "vj":
		return Credential{
			User:      "viajor",
			Passwd:    "3HO@ooX-GW",
			BranchID:  "1481",
			AgencyID:  "1874",
			AgentSign: "VJ",
		}
	case user == "yz":
		return Credential{
			User:      "YZZERGW",
			Passwd:    "Mh3cYpRR2f",
			BranchID:  "256848",
			AgencyID:  "256848",
			AgentSign: "YZ",
		}
	case user == "gl":
		return Credential{
			User:     "globalis",
			Passwd:   "GlB3cYpRc2f",
			BranchID: "1000",
			AgencyID: "1000",
		}
	case user == "mc":
		return Credential{
			User:     "ZRNGW",
			Passwd:   "42JUQDfE",
			BranchID: "213",
			AgencyID: "2900",
		}
	case user == "htl":
		return Credential{
			User:     "hotelli",
			Passwd:   "s3u*ukA#rAnatruC",
			BranchID: "1000",
			AgencyID: "254984",
		}
	case user == "avt":
		return Credential{
			User:     "avantrip",
			Passwd:   "bibamgroup@2018",
			BranchID: "1000",
		}
	case user == "kw":
		return Credential{
			User:     "KIWIGWINT",
			Passwd:   "integracaogw",
			BranchID: "1000",
		}
	case user == "is":
		return Credential{
			User:     "Inspireegw",
			Passwd:   "integracaogw",
			BranchID: "1000",
		}
	case user == "it":
		return Credential{
			User:     "itravel",
			Passwd:   "itravel31541000",
			BranchID: "13139",
			AgencyID: "13139",
		}
	case user == "wo":
		return Credential{
			User:     "woobagw",
			Passwd:   "integracaogw",
			BranchID: "1000",
		}

	case user == "in":
		return Credential{
			User:     "infarecvc",
			Passwd:   "2WeR61oy9I4m",
			BranchID: "1000",
		}

	case user == "md":
		return Credential{
			User:     "md_afl",
			Passwd:   "45Jr1ze2",
			BranchID: "1020",
		}

	case user == "pt":
		return Credential{
			User:     "proponto",
			Passwd:   "4nmd7Yh4",
			BranchID: "1008",
		}

	case user == "wb":
		return Credential{
			User:     "wooba",
			Passwd:   "woobagw",
			BranchID: "1000",
		}

	case user == "es":
		return Credential{
			User:     "esfera",
			Passwd:   "EKoMD9am",
			BranchID: "1035",
		}

	case user == "trafega":
		return Credential{
			User:     "trafega",
			Passwd:   "MhmcYpMM16e",
			BranchID: "1000",
		}

	}

	return Credential{}
}

func GetURL(host string) string {
	//dev|mk|cvc-hom|cvc-dev|cvc|lvl|lojas|int|sv|ra|mm|orders
	switch {
	case host == "cs":
		return "http://localhost:8500/v1/catalog/service"
	case host == "dev":
		return "http://localhost:8080"
		//		return "http://localhost:8082"
		//return "http://localhost:9999"
	case host == "mk":
		return "http://192.168.99.100:31020"
	case host == "cvc-hom":
		//return "https://search-cvc-hom.reservafacil.tur.br"
		return "https://gwa-cvc-hom.reservafacil.tur.br"
	case host == "cvc-dev":
		return "https://search-cvc-dev.reservafacil.tur.br"
	case host == "pre":
		return "https://us-search-pre-prod.reservafacil.tur.br"
	case host == "cvc":
		return "https://internal-search-cvc-prod.reservafacil.tur.br"
	case host == "lvl":
		return "https://internal-search-lvl-prod.reservafacil.tur.br"
	case host == "lojas":
		return "https://internal-search-lojas-prod.reservafacil.tur.br"
	case host == "int":
		return "https://search-int-prod.reservafacil.tur.br"
	case host == "sv":
		return "https://internal-search-sv-prod.reservafacil.tur.br"
	case host == "ra":
		return "https://internal-search-ra-prod.reservafacil.tur.br"
	case host == "afl":
		return "https://internal-search-afl-prod.reservafacil.tur.br"
	case host == "afl-us":
		return "https://internal-us-search-afl-prod.reservafacil.tur.br"
	case host == "mm":
		return "https://us-search-mm-prod.reservafacil.tur.br"
	case host == "vs":
		return "https://search-vis-prod.reservafacil.tur.br"
	case host == "tr":
		return "https://search-trd-prod.reservafacil.tur.br"
	case host == "orders":
		return "https://orders-cvc-prod.reservafacil.tur.br"
	case host == "afl-gcp":
		return "https://search-afl-poc.reservafacil.tur.br"
	case host == "pre":
		return "https://us-search-pre-prod.reservafacil.tur.br"

	}
	return host
}

func GetURLPrecificador(host string) string {
	//dev|mk|cvc-hom|cvc-dev|cvc|lvl|lojas|int|sv|ra|mm|orders
	switch {
	case host == "dev":
		//		return "http://localhost:8480/precificador"
		return "https://internal-elb-precificador-cvc-prod.reservafacil.tur.br/precificador"
	case host == "mk":
		return "http://192.168.99.100:31020/precificador"
	case host == "cvc-hom":
		return "https://internal-elb-precificador-cvc-hom.reservafacil.tur.br/precificador"
	case host == "cvc-dev":
		return "https://internal-elb-precificador-cvc-dev.reservafacil.tur.br/precificador"
	default:
		return "https://internal-alb-precific-cvc-prod.reservafacil.tur.br/precificador"
	}
}

func GetURLGwapi(host string) string {
	//dev|mk|cvc-hom|cvc-dev|cvc|lvl|lojas|int|sv|ra|mm|orders
	switch {
	case host == "dev":
		return "http://localhost:8480/gwapi"
	case host == "mk":
		return "http://192.168.99.100:31029/gwapi"
	case host == "cvc-hom":
		return "https://internal-elb-gwapi-cvc-hom.reservafacil.tur.br/gwapi"
	case host == "cvc-dev":
		return "https://internal-elb-gwapi-cvc-dev.reservafacil.tur.br/gwapi"
	default:
		return "https://internal-elb-gwapi-cvc-prod.reservafacil.tur.br/gwapi"
	}
}

func GetURLNotify(host string) string {
	//dev|mk|cvc-hom|cvc-dev|cvc|lvl|lojas|int|sv|ra|mm|orders
	switch {
	case host == "dev":
		return "http://localhost:8080/gwaereo"
	case host == "cvc-hom":
		return "https://todo"
	case host == "cvc-dev":
		return "https://3t26omfcp4.execute-api.us-east-1.amazonaws.com/v1"
	default:
		return "https://3t26omfcp4.execute-api.us-east-1.amazonaws.com/v1"
	}
}

func GetJWTPrec(urlPrecificador string, username string, password string, branchId string, agencyId string) string {
	req, err := http.NewRequest("GET", urlPrecificador+"/auth/v1/autentica/login", nil)
	req.Header.Add("X-AUTH-LOGIN", username)
	req.Header.Add("X-AUTH-PASS", password)
	req.Header.Add("X-AUTH-BRANCH", branchId)
	req.Header.Add("X-AUTH-AGENCY", agencyId)

	log.Trace(req)

	client := http.DefaultClient
	resp, err := client.Do(req)

	log.Trace(resp)

	if err != nil {
		log.Info("Error", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	response := JWT{}
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Info("Error", err)
	}

	return response.Value
}

type JWT struct {
	Value string `json:"jwt"`
}

type Date string

func (s Date) ShortDate() string {
	t, err := time.Parse(YYYY_MM_DD_T_HH_MM_SS, string(s))
	if err != nil {
		log.Info("Error parsing date", s, err)
		os.Exit(1)
	}
	return t.Format(SHORT_DATE)
}

func (s Date) ShortTime() string {
	t, err := time.Parse(YYYY_MM_DD_T_HH_MM_SS, string(s))
	if err != nil {
		log.Info("Error parsing date", s, err)
		os.Exit(1)
	}
	return t.Format(SHORT_TIME)
}

func GetSource(sourceKey string) string {
	switch {
	case "LAT" == sourceKey:
		return "LATAM"
	case "AVI" == sourceKey:
		return "AVIANCA"
	case "GOL" == sourceKey:
		return "GOL"
	case "AZUL" == sourceKey:
		return "AZUL"
	case "AMD" == sourceKey:
		return "AMADEUS"
	case "SAB" == sourceKey:
		return "SABRE"
	case "FLX" == sourceKey:
		return "FARELOGIX"
	case "FRT" == sourceKey:
		return "FRETAMENTO"
	case "CION" == sourceKey:
		return "CIONS"
	}
	return ""
}
